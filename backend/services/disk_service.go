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

type DiskService interface {
	StartCollecting()
	StopCollecting()
	GetLatestLatency(limit int) ([]models.DiskLatency, error)
}

type diskService struct {
	repo      repository.DiskRepository
	collector *collector.DiskCollector
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewDiskService(repo repository.DiskRepository) DiskService {
	ctx, cancel := context.WithCancel(context.Background())
	return &diskService{
		repo:      repo,
		collector: collector.NewDiskCollector(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *diskService) StartCollecting() {
	// Start the streaming collector
	if err := s.collector.Start(); err != nil {
		log.Printf("Failed to start disk collector: %v", err)
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		log.Println("Disk latency collector started")

		for {
			select {
			case <-s.ctx.Done():
				log.Println("Disk latency collector stopped")
				return
			case <-ticker.C:
				latencies := s.collector.GetEvents()
				if len(latencies) > 0 {
					if err := s.repo.SaveLatencySnapshot(latencies); err != nil {
						log.Printf("Error saving disk latency: %v", err)
					}
				}
			}
		}
	}()
}

func (s *diskService) StopCollecting() {
	s.cancel()
	s.collector.Stop()
	s.wg.Wait()
}

func (s *diskService) GetLatestLatency(limit int) ([]models.DiskLatency, error) {
	return s.repo.GetLatestLatency(limit)
}
