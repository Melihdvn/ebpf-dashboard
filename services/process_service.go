package services

import (
	"ebpf-dashboard/collector"
	"ebpf-dashboard/models"
	"ebpf-dashboard/repository"
	"log"
	"sync"
	"time"
)

type ProcessService interface {
	StartCollecting()
	StopCollecting()
	GetRecentProcesses(limit int) ([]models.ProcessEvent, error)
}

type processService struct {
	repo      repository.ProcessRepository
	collector *collector.ProcessCollector
	stopChan  chan bool
	wg        sync.WaitGroup
}

func NewProcessService(repo repository.ProcessRepository) ProcessService {
	return &processService{
		repo:      repo,
		collector: collector.NewProcessCollector(),
		stopChan:  make(chan bool),
	}
}

func (s *processService) StartCollecting() {
	// Start the continuous collector
	if err := s.collector.Start(); err != nil {
		log.Printf("Failed to start process collector: %v", err)
		return
	}

	// Start goroutine to periodically save events
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopChan:
				s.collector.Stop()
				return
			case <-ticker.C:
				// Get accumulated events
				events := s.collector.GetEvents()

				// Save them to database
				for _, event := range events {
					if err := s.repo.SaveProcess(event); err != nil {
						log.Printf("Error saving process: %v", err)
					}
				}
			}
		}
	}()
}

func (s *processService) StopCollecting() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *processService) GetRecentProcesses(limit int) ([]models.ProcessEvent, error) {
	return s.repo.GetRecentProcesses(limit)
}
