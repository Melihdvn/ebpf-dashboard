package models

import "time"

// TCPLifeEvent represents a TCP connection lifecycle event
type TCPLifeEvent struct {
	ID         int       `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	PID        int       `json:"pid"`
	Comm       string    `json:"comm"`
	LocalAddr  string    `json:"local_addr"`
	LocalPort  int       `json:"local_port"`
	RemoteAddr string    `json:"remote_addr"`
	RemotePort int       `json:"remote_port"`
	TxKB       float64   `json:"tx_kb"`
	RxKB       float64   `json:"rx_kb"`
	DurationMS float64   `json:"duration_ms"`
}
