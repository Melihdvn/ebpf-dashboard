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

type NetworkService interface {
	StartCollecting()
	StopCollecting()
	GetRecentConnections(limit int) ([]models.NetworkConnection, error)
}

type networkService struct {
	repo      repository.NetworkRepository
	collector *collector.NetworkCollector
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewNetworkService(repo repository.NetworkRepository) NetworkService {
	ctx, cancel := context.WithCancel(context.Background())
	return &networkService{
		repo:      repo,
		collector: collector.NewNetworkCollector(),
		ctx:       ctx,
		cancel:    cancel,
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
			case <-s.ctx.Done():
				s.collector.Stop()
				return
			case <-ticker.C:
				// Get accumulated events
				events := s.collector.GetEvents()

				// Save them to database using batch insert
				if len(events) > 0 {
					if err := s.repo.SaveConnections(events); err != nil {
						log.Printf("Error saving connections: %v", err)
					}
				}
			}
		}
	}()
}

func (s *networkService) StopCollecting() {
	s.cancel()
	s.wg.Wait()
}

func (s *networkService) GetRecentConnections(limit int) ([]models.NetworkConnection, error) {
	return s.repo.GetRecentConnections(limit)
}
