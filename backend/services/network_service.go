package services

import (
	"ebpf-dashboard/collector"
	"ebpf-dashboard/models"
	"ebpf-dashboard/repository"
	"log"
	"sync"
	"time"
)

type NetworkService interface {
	StartCollecting()
	StopCollecting()
	GetRecentConnections(limit int) ([]models.NetworkConnection, error)
}

type networkService struct {
	repo      repository.NetworkRepository
	collector *collector.NetworkCollector
	stopChan  chan bool
	wg        sync.WaitGroup
}

func NewNetworkService(repo repository.NetworkRepository) NetworkService {
	return &networkService{
		repo:      repo,
		collector: collector.NewNetworkCollector(),
		stopChan:  make(chan bool),
	}
}

func (s *networkService) StartCollecting() {
	// Start the continuous collector
	if err := s.collector.Start(); err != nil {
		log.Printf("Failed to start network collector: %v", err)
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
					if err := s.repo.SaveConnection(event); err != nil {
						log.Printf("Error saving connection: %v", err)
					}
				}
			}
		}
	}()
}

func (s *networkService) StopCollecting() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *networkService) GetRecentConnections(limit int) ([]models.NetworkConnection, error) {
	return s.repo.GetRecentConnections(limit)
}
