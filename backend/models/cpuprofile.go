package models

import "time"

// CPUProfile represents a CPU profiling sample with stack trace
type CPUProfile struct {
	ID          int       `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessName string    `json:"process_name"`
	StackTrace  string    `json:"stack_trace"`
	SampleCount int       `json:"sample_count"`
}
