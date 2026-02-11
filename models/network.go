package models

import "time"

type NetworkConnection struct {
	ID         int       `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	PID        string    `json:"pid"`
	Comm       string    `json:"comm"`
	IPVersion  string    `json:"ip_version"`
	SourceAddr string    `json:"source_addr"`
	SourcePort string    `json:"source_port"`
	DestAddr   string    `json:"dest_addr"`
	DestPort   string    `json:"dest_port"`
}
