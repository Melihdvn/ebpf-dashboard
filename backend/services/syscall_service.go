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

type SyscallService interface {
	Start()
	Stop()
	GetRecentStats(limit int) ([]models.SyscallStat, error)
}

type syscallService struct {
	repo      *repository.SyscallRepository
	collector *collector.SyscallCollector
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewSyscallService(repo *repository.SyscallRepository) SyscallService {
	ctx, cancel := context.WithCancel(context.Background())
	return &syscallService{
		repo:      repo,
		collector: collector.NewSyscallCollector(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start begins periodic syscall statistics collection
func (s *syscallService) Start() {
	// Start the streaming collector
	if err := s.collector.Start(); err != nil {
		log.Printf("Failed to start syscall collector: %v", err)
		return
	}

	s.wg.Add(1)
	go s.collectPeriodically()
	log.Println("Syscall stats service started")
}

// Stop gracefully stops the syscall stats service
func (s *syscallService) Stop() {
	s.cancel()
	s.collector.Stop()
	s.wg.Wait()
	log.Println("Syscall stats service stopped")
}

func (s *syscallService) collectPeriodically() {
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

func (s *syscallService) collectAndSave() {
	stats := s.collector.GetEvents()

	if len(stats) == 0 {
		return
	}

	if err := s.repo.SaveSyscallStats(stats); err != nil {
		log.Printf("Error saving syscall stats: %v", err)
		return
	}
}

// GetRecentStats retrieves recent syscall statistics
func (s *syscallService) GetRecentStats(limit int) ([]models.SyscallStat, error) {
	return s.repo.GetRecentSyscallStats(limit)
}
