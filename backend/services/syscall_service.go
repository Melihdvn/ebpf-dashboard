package services

import (
	"context"
	"ebpf-dashboard/collector"
	"ebpf-dashboard/repository"
	"log"
	"sync"
	"time"
)

type SyscallService struct {
	repo   *repository.SyscallRepository
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewSyscallService(repo *repository.SyscallRepository) *SyscallService {
	ctx, cancel := context.WithCancel(context.Background())
	return &SyscallService{
		repo:   repo,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins periodic syscall statistics collection
func (s *SyscallService) Start() {
	s.wg.Add(1)
	go s.collectPeriodically()
	log.Println("Syscall stats service started")
}

// Stop gracefully stops the syscall stats service
func (s *SyscallService) Stop() {
	s.cancel()
	s.wg.Wait()
	log.Println("Syscall stats service stopped")
}

func (s *SyscallService) collectPeriodically() {
	defer s.wg.Done()

	// Collect immediately on start
	s.collectAndSave()

	// Then collect every 5 seconds (syscount runs for 5 seconds, so essentially continuous)
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

func (s *SyscallService) collectAndSave() {
	stats, err := collector.CollectSyscallStats()
	if err != nil {
		log.Printf("Error collecting syscall stats: %v", err)
		return
	}

	if len(stats) == 0 {
		return
	}

	if err := s.repo.SaveSyscallStats(stats); err != nil {
		log.Printf("Error saving syscall stats: %v", err)
		return
	}
}

// GetRecentStats retrieves recent syscall statistics
func (s *SyscallService) GetRecentStats(limit int) (interface{}, error) {
	return s.repo.GetRecentSyscallStats(limit)
}
