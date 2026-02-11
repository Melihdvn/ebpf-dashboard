package models

import "time"

type ProcessEvent struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Time      string    `json:"time"`
	PID       string    `json:"pid"`
	Comm      string    `json:"comm"`
	Args      string    `json:"args"`
}
