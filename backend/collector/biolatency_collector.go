package collector

import (
	"bufio"
	"context"
	"ebpf-dashboard/models"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type DiskCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.DiskLatency
	mu      sync.Mutex
	running bool
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		events: make(chan models.DiskLatency, 100),
	}
}

func (c *DiskCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start biolatency in continuous mode: 1 second intervals
	c.cmd = exec.CommandContext(ctx, "sudo", "biolatency", "1")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		log.Printf("Failed to start biolatency: %v", err)
		return err
	}

	c.running = true
	log.Println("biolatency collector started")

	// Read output in a goroutine
	go func() {
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			log.Println("biolatency collector stopped")
		}()

		reader := bufio.NewReader(stdout)
		re := regexp.MustCompile(`(\d+)\s*->\s*(\d+)\s*:\s*(\d+)`)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("biolatency read error: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Parse histogram lines
			matches := re.FindStringSubmatch(line)
			if len(matches) == 4 {
				rangeMin, _ := strconv.Atoi(matches[1])
				rangeMax, _ := strconv.Atoi(matches[2])
				count, _ := strconv.Atoi(matches[3])

				latency := models.DiskLatency{
					RangeMin: rangeMin,
					RangeMax: rangeMax,
					Count:    count,
				}

				// Send to channel (non-blocking)
				select {
				case c.events <- latency:
				default:
					// Channel full, skip this event
				}
			}
		}
	}()

	return nil
}

func (c *DiskCollector) Stop() {
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

func (c *DiskCollector) GetEvents() []models.DiskLatency {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []models.DiskLatency

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
