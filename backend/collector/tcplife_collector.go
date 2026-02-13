package collector

import (
	"bufio"
	"context"
	"ebpf-dashboard/models"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type TCPLifeCollector struct {
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	events  chan models.TCPLifeEvent
	mu      sync.Mutex
	running bool
}

func NewTCPLifeCollector() *TCPLifeCollector {
	return &TCPLifeCollector{
		events: make(chan models.TCPLifeEvent, 100),
	}
}

func (c *TCPLifeCollector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// Start tcplife in continuous mode
	// Use stdbuf to disable output buffering
	// Command: stdbuf -oL tcplife
	c.cmd = exec.CommandContext(ctx, "stdbuf", "-oL", "tcplife")

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		log.Printf("Failed to start tcplife: %v", err)
		return err
	}

	c.running = true
	log.Println("tcplife collector started")

	// Read output line by line in a goroutine
	go func() {
		defer func() {
			c.mu.Lock()
			c.running = false
			c.mu.Unlock()
			log.Println("tcplife collector stopped")
		}()

		scanner := bufio.NewScanner(stdout)

		// Skip header line
		// PID     COMM             LADDR           LPORT RADDR           RPORT TX_KB  RX_KB  MS

		for scanner.Scan() {
			line := scanner.Text()

			// Skip header and empty lines
			if strings.HasPrefix(line, "PID") || strings.TrimSpace(line) == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 9 {
				// Parse fields
				pid, _ := strconv.Atoi(fields[0])
				comm := fields[1]

				// Handle COMM field that might have spaces or be "Socket Thread"
				// Determining where COMM ends and LADDR starts is tricky with spaces
				// But we know LADDR is an IP address

				// Let's try to parse from the end since the last columns are fixed
				// 9 columns minimum.
				// Last 3 are metrics: TX_KB, RX_KB, MS
				// Then RPORT, RADDR, LPORT, LADDR
				// The rest at the beginning are PID and COMM

				n := len(fields)
				if n < 9 {
					continue
				}

				durationMS, _ := strconv.ParseFloat(fields[n-1], 64)
				rxKB, _ := strconv.ParseFloat(fields[n-2], 64)
				txKB, _ := strconv.ParseFloat(fields[n-3], 64)
				remotePort, _ := strconv.Atoi(fields[n-4])
				remoteAddr := fields[n-5]
				localPort, _ := strconv.Atoi(fields[n-6])
				localAddr := fields[n-7]

				// Reconstruct COMM from fields[1] to fields[n-8]
				commParts := fields[1 : n-7]
				comm = strings.Join(commParts, " ")

				event := models.TCPLifeEvent{
					PID:        pid,
					Comm:       comm,
					LocalAddr:  localAddr,
					LocalPort:  localPort,
					RemoteAddr: remoteAddr,
					RemotePort: remotePort,
					TxKB:       txKB,
					RxKB:       rxKB,
					DurationMS: durationMS,
				}

				// Send to channel (non-blocking)
				select {
				case c.events <- event:
				default:
					// Channel full, skip this event
				}
			}
		}
	}()

	return nil
}

func (c *TCPLifeCollector) Stop() {
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

func (c *TCPLifeCollector) GetEvents() []models.TCPLifeEvent {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []models.TCPLifeEvent

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
