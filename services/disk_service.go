package services

import (
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
	repo     repository.DiskRepository
	stopChan chan bool
	wg       sync.WaitGroup
}

func NewDiskService(repo repository.DiskRepository) DiskService {
	return &diskService{
		repo:     repo,
		stopChan: make(chan bool),
	}
}

func (s *diskService) StartCollecting() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		log.Println("Disk latency collector started")

		for {
			select {
			case <-s.stopChan:
				log.Println("Disk latency collector stopped")
				return
			case <-ticker.C:
				latencies, err := collector.CollectDiskLatency()
				if err != nil {
					log.Printf("Error collecting disk latency: %v", err)
					continue
				}

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
	close(s.stopChan)
	s.wg.Wait()
}

func (s *diskService) GetLatestLatency(limit int) ([]models.DiskLatency, error) {
	return s.repo.GetLatestLatency(limit)
}
