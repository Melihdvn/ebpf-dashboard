package services

import (
	"context"
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
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewProcessService(repo repository.ProcessRepository) ProcessService {
	ctx, cancel := context.WithCancel(context.Background())
	return &processService{
		repo:      repo,
		collector: collector.NewProcessCollector(),
		ctx:       ctx,
		cancel:    cancel,
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
			case <-s.ctx.Done():
				s.collector.Stop()
				return
			case <-ticker.C:
				// Get accumulated events
				events := s.collector.GetEvents()

				// Save them to database using batch insert
				if len(events) > 0 {
					if err := s.repo.SaveProcesses(events); err != nil {
						log.Printf("Error saving processes: %v", err)
					}
				}
			}
		}
	}()
}

func (s *processService) StopCollecting() {
	s.cancel()
	s.wg.Wait()
}

func (s *processService) GetRecentProcesses(limit int) ([]models.ProcessEvent, error) {
	return s.repo.GetRecentProcesses(limit)
}
