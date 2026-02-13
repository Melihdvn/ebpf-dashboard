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

type TCPLifeService interface {
	StartCollecting() error
	StopCollecting()
	GetRecentEvents(limit int) ([]models.TCPLifeEvent, error)
}

type tcpLifeService struct {
	collector *collector.TCPLifeCollector
	repo      *repository.TCPLifeRepository
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewTCPLifeService(repo *repository.TCPLifeRepository) TCPLifeService {
	ctx, cancel := context.WithCancel(context.Background())
	return &tcpLifeService{
		collector: collector.NewTCPLifeCollector(),
		repo:      repo,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// StartCollecting starts the background collection process
func (s *tcpLifeService) StartCollecting() error {
	if err := s.collector.Start(); err != nil {
		return err
	}

	s.wg.Add(1)
	go s.processEvents()

	return nil
}

// StopCollecting stops the background collection process
func (s *tcpLifeService) StopCollecting() {
	s.cancel()
	s.collector.Stop()
	s.wg.Wait()
}

func (s *tcpLifeService) processEvents() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
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
func (s *tcpLifeService) GetRecentEvents(limit int) ([]models.TCPLifeEvent, error) {
	return s.repo.GetRecentTCPLifeEvents(limit)
}
