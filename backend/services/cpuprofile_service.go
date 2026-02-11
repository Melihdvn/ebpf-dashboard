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

type CPUProfileService interface {
	Start()
	Stop()
	GetRecentProfiles(limit int) ([]models.CPUProfile, error)
}

type cpuProfileService struct {
	repo      *repository.CPUProfileRepository
	collector *collector.CPUProfileCollector
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewCPUProfileService(repo *repository.CPUProfileRepository) CPUProfileService {
	ctx, cancel := context.WithCancel(context.Background())
	return &cpuProfileService{
		repo:      repo,
		collector: collector.NewCPUProfileCollector(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start begins periodic CPU profiling collection
func (s *cpuProfileService) Start() {
	// Start the streaming collector
	if err := s.collector.Start(); err != nil {
		log.Printf("Failed to start CPU profile collector: %v", err)
		return
	}

	s.wg.Add(1)
	go s.collectPeriodically()
	log.Println("CPU profile service started")
}

// Stop gracefully stops the CPU profiling service
func (s *cpuProfileService) Stop() {
	s.cancel()
	s.collector.Stop()
	s.wg.Wait()
	log.Println("CPU profile service stopped")
}

func (s *cpuProfileService) collectPeriodically() {
	defer s.wg.Done()

	// Collect every 5 seconds
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

func (s *cpuProfileService) collectAndSave() {
	profiles := s.collector.GetEvents()

	if len(profiles) == 0 {
		return
	}

	if err := s.repo.SaveCPUProfiles(profiles); err != nil {
		log.Printf("Error saving CPU profiles: %v", err)
		return
	}

	log.Printf("Saved %d CPU profile samples", len(profiles))
}

// GetRecentProfiles retrieves recent CPU profile data
func (s *cpuProfileService) GetRecentProfiles(limit int) ([]models.CPUProfile, error) {
	return s.repo.GetRecentCPUProfiles(limit)
}
