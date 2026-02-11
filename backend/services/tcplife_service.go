package services

import (
	"ebpf-dashboard/collector"
	"ebpf-dashboard/repository"
	"log"
	"sync"
	"time"
)

type TCPLifeService struct {
	collector *collector.TCPLifeCollector
	repo      *repository.TCPLifeRepository
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

func NewTCPLifeService(repo *repository.TCPLifeRepository) *TCPLifeService {
	return &TCPLifeService{
		collector: collector.NewTCPLifeCollector(),
		repo:      repo,
		stopChan:  make(chan struct{}),
	}
}

// StartCollecting starts the background collection process
func (s *TCPLifeService) StartCollecting() error {
	if err := s.collector.Start(); err != nil {
		return err
	}

	s.wg.Add(1)
	go s.processEvents()

	return nil
}

// StopCollecting stops the background collection process
func (s *TCPLifeService) StopCollecting() {
	s.collector.Stop()
	close(s.stopChan)
	s.wg.Wait()
}

func (s *TCPLifeService) processEvents() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			events := s.collector.GetEvents()
			if len(events) > 0 {
				if err := s.repo.SaveTCPLifeEvents(events); err != nil {
					log.Printf("Error saving tcplife events: %v", err)
				}
			}
		}
	}
}

// GetRecentEvents retrieves recent TCP lifecycle events
func (s *TCPLifeService) GetRecentEvents(limit int) (interface{}, error) {
	return s.repo.GetRecentTCPLifeEvents(limit)
}
