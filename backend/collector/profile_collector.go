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

type CPUProfileCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.CPUProfile
	mu      sync.Mutex
	running bool
}

func NewCPUProfileCollector() *CPUProfileCollector {
	return &CPUProfileCollector{
		events: make(chan models.CPUProfile, 100),
	}
}

func (c *CPUProfileCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start profile-bpfcc: sample at 99 Hz, continuous mode with 5 second intervals
	c.cmd = exec.CommandContext(ctx, "sudo", "profile-bpfcc", "-F", "99", "5")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		log.Printf("Failed to start profile-bpfcc: %v", err)
		return err
	}

	c.running = true
	log.Println("profile-bpfcc collector started")

	// Read output in a goroutine
	go func() {
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			log.Println("profile-bpfcc collector stopped")
		}()

		reader := bufio.NewReader(stdout)
		var currentStack []string
		var processName string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("profile-bpfcc read error: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)

			// Skip header and empty lines
			if line == "" || strings.HasPrefix(line, "Sampling") {
				// If we have accumulated a stack, save it
				if len(currentStack) > 0 && processName != "" {
					stackTrace := strings.Join(currentStack, "\n")
					profile := models.CPUProfile{
						ProcessName: processName,
						StackTrace:  stackTrace,
						SampleCount: 1, // Will be aggregated in service
					}

					// Send to channel (non-blocking)
					select {
					case c.events <- profile:
					default:
						// Channel full, skip this event
					}

					// Reset for next stack
					currentStack = []string{}
					processName = ""
				}
				continue
			}

			// Check if this line is a number (sample count)
			if matched, _ := regexp.MatchString(`^\d+$`, line); matched {
				// This is a count, process the accumulated stack
				if len(currentStack) > 0 {
					// Last item in stack is the process name
					processName = currentStack[len(currentStack)-1]
					// Remove process name from stack
					stackLines := currentStack[:len(currentStack)-1]
					stackTrace := strings.Join(stackLines, "\n")

					count, _ := strconv.Atoi(line)

					profile := models.CPUProfile{
						ProcessName: processName,
						StackTrace:  stackTrace,
						SampleCount: count,
					}

					// Send to channel (non-blocking)
					select {
					case c.events <- profile:
					default:
						// Channel full, skip this event
					}

					// Reset for next stack
					currentStack = []string{}
					processName = ""
				}
				continue
			}

			// This is a stack frame, add to current stack
			currentStack = append(currentStack, line)
		}
	}()

	return nil
}

func (c *CPUProfileCollector) Stop() {
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

func (c *CPUProfileCollector) GetEvents() []models.CPUProfile {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []models.CPUProfile

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
