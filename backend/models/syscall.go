package models

import "time"

// SyscallStat represents system call statistics
type SyscallStat struct {
	ID          int       `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	SyscallName string    `json:"syscall_name"`
	Count       int       `json:"count"`
}
