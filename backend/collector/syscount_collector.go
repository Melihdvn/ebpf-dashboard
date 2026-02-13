package collector

import (
	"bufio"
	"context"
	"ebpf-dashboard/models"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type SyscallCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.SyscallStat
	mu      sync.Mutex
	running bool
}

func NewSyscallCollector() *SyscallCollector {
	return &SyscallCollector{
		events: make(chan models.SyscallStat, 100),
	}
}

func (c *SyscallCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start syscount-bpfcc: 5 second intervals, continuous mode
	c.cmd = exec.CommandContext(ctx, "sudo", "syscount-bpfcc", "-i", "5")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		log.Printf("Failed to start syscount-bpfcc: %v", err)
		return err
	}

	c.running = true
	log.Println("syscount-bpfcc collector started")

	// Read output in a goroutine
	go func() {
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			log.Println("syscount-bpfcc collector stopped")
		}()

		reader := bufio.NewReader(stdout)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("syscount-bpfcc read error: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)

			// Skip empty lines, headers, and tracing messages
			if line == "" || strings.HasPrefix(line, "SYSCALL") ||
				strings.Contains(line, "Tracing") || strings.HasPrefix(line, "[") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 2 {
				syscallName := fields[0]
				countStr := fields[1]

				count, err := strconv.Atoi(countStr)
				if err == nil {
					stat := models.SyscallStat{
						SyscallName: syscallName,
						Count:       count,
					}

					// Send to channel (non-blocking)
					select {
					case c.events <- stat:
					default:
						// Channel full, skip this event
					}
				}
			}
		}
	}()

	return nil
}

func (c *SyscallCollector) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	if c.cancel != nil {
		c.cancel()
	}

	if c.cmd != nil && c.cmd.Process != nil {
		c.cmd.Process.Kill()
	}

	c.running = false
}

func (c *SyscallCollector) GetEvents() []models.SyscallStat {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []models.SyscallStat

	// Drain the channel
	for {
		select {
		case event := <-c.events:
			events = append(events, event)
		default:
			return events
		}
	}
}
