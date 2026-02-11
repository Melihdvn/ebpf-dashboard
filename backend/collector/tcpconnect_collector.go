package collector

import (
	"bufio"
	"context"
	"ebpf-dashboard/models"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type NetworkCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.NetworkConnection
	mu      sync.Mutex
	running bool
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{
		events: make(chan models.NetworkConnection, 100),
	}
}

func (c *NetworkCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start tcpconnect in continuous mode (no sudo needed, app runs with sudo)
	// Use stdbuf to disable output buffering so we get lines immediately
	c.cmd = exec.CommandContext(ctx, "stdbuf", "-oL", "tcpconnect")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		log.Printf("Failed to start tcpconnect: %v", err)
		return err
	}

	c.running = true
	log.Println("tcpconnect collector started")

	// Read output line by line in a goroutine
	go func() {
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			log.Println("tcpconnect collector stopped")
		}()

		reader := bufio.NewReader(stdout)
		lineNum := 0

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("tcpconnect read error: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			lineNum++

			// Skip header line and empty lines
			if lineNum == 1 || line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 6 {
				// tcpconnect output format: PID COMM IP SADDR DADDR DPORT
				conn := models.NetworkConnection{
					PID:        fields[0],
					Comm:       fields[1],
					IPVersion:  "IPv" + fields[2],
					SourceAddr: fields[3],
					DestAddr:   fields[4],
					DestPort:   fields[5],
				}

				// Send to channel (non-blocking)
				select {
				case c.events <- conn:
				default:
					// Channel full, skip this event
				}
			}
		}
	}()

	return nil
}

func (c *NetworkCollector) Stop() {
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

func (c *NetworkCollector) GetEvents() []models.NetworkConnection {
	var events []models.NetworkConnection

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
