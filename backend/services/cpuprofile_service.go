package services

import (
	"context"
	"ebpf-dashboard/collector"
	"ebpf-dashboard/repository"
	"log"
	"sync"
	"time"
)

type CPUProfileService struct {
	repo   *repository.CPUProfileRepository
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewCPUProfileService(repo *repository.CPUProfileRepository) *CPUProfileService {
	ctx, cancel := context.WithCancel(context.Background())
	return &CPUProfileService{
		repo:   repo,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins periodic CPU profiling collection
func (s *CPUProfileService) Start() {
	s.wg.Add(1)
	go s.collectPeriodically()
	log.Println("CPU profile service started")
}

// Stop gracefully stops the CPU profiling service
func (s *CPUProfileService) Stop() {
	s.cancel()
	s.wg.Wait()
	log.Println("CPU profile service stopped")
}

func (s *CPUProfileService) collectPeriodically() {
	defer s.wg.Done()

	// Collect immediately on start
	s.collectAndSave()

	// Then collect every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.collectAndSave()
		}
	}
}

func (s *CPUProfileService) collectAndSave() {
	profiles, err := collector.CollectCPUProfile()
	if err != nil {
		log.Printf("Error collecting CPU profiles: %v", err)
		return
	}

	if len(profiles) == 0 {
		log.Println("No CPU profile data collected")
		return
	}

	if err := s.repo.SaveCPUProfiles(profiles); err != nil {
		log.Printf("Error saving CPU profiles: %v", err)
		return
	}

	log.Printf("Saved %d CPU profile samples", len(profiles))
}

// GetRecentProfiles retrieves recent CPU profile data
func (s *CPUProfileService) GetRecentProfiles(limit int) (interface{}, error) {
	return s.repo.GetRecentCPUProfiles(limit)
}
