package collector

import (
	"ebpf-dashboard/models"
	"ebpf-dashboard/utils"
	"strconv"
	"strings"
	"time"
)

// CollectSyscallStats runs syscount-bpfcc and returns parsed syscall statistics
func CollectSyscallStats() ([]models.SyscallStat, error) {
	// Run syscount: interval 5 seconds, count 1
	// Command: syscount-bpfcc -i 5 1
	output, err := utils.ExecuteWithSudo(10*time.Second, "syscount-bpfcc", "-i", "5", "1")
	if err != nil {
		// Even with error, we might have partial output
		if len(output) == 0 {
			return nil, err
		}
	}

	return parseSyscountOutput(string(output)), nil
}

func parseSyscountOutput(output string) []models.SyscallStat {
	// syscount output format:
	// SYSCALL                   COUNT
	// read                       2245
	// futex                      1487

	var stats []models.SyscallStat
	lines := strings.Split(output, "\n")

	// Skip header lines (usually first 2-3 lines depending on output)
	// We'll look for lines that have 2 fields: string and int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "SYSCALL") || strings.Contains(line, "Tracing") || strings.HasPrefix(line, "[") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			syscallName := fields[0]
			countStr := fields[1]

			count, err := strconv.Atoi(countStr)
			if err == nil {
				stat := models.SyscallStat{
					SyscallName: syscallName,
					Count:       count,
				}
				stats = append(stats, stat)
			}
		}
	}

	return stats
}
