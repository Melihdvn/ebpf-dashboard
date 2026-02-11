package collector

import (
	"bufio"
	"context"
	"ebpf-dashboard/models"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type ProcessCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.ProcessEvent
	mu      sync.Mutex
	running bool
}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{
		events: make(chan models.ProcessEvent, 100),
	}
}

func (c *ProcessCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start execsnoop in continuous mode (no sudo needed, app runs with sudo)
	c.cmd = exec.CommandContext(ctx, "execsnoop", "-T")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		return err
	}

	c.running = true
	log.Println("execsnoop collector started")

	// Read output line by line in a goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		lineNum := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineNum++

			// Skip header line
			if lineNum == 1 || strings.TrimSpace(line) == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 5 {
				event := models.ProcessEvent{
					Time: fields[0],
					PID:  fields[2],
					Comm: fields[1],
					Args: strings.Join(fields[5:], " "),
				}

				// Send to channel (non-blocking)
				select {
				case c.events <- event:
				default:
					// Channel full, skip this event
				}
			}
		}

		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
		log.Println("execsnoop collector stopped")
	}()

	return nil
}

func (c *ProcessCollector) Stop() {
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

func (c *ProcessCollector) GetEvents() []models.ProcessEvent {
	var events []models.ProcessEvent

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
