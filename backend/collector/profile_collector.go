package collector

import (
	"ebpf-dashboard/models"
	"ebpf-dashboard/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CollectCPUProfile runs profile-bpfcc and returns parsed CPU profiling data
func CollectCPUProfile() ([]models.CPUProfile, error) {
	// Run profile-bpfcc: sample at 99 Hz for 1 second
	// -F 99: frequency (99 samples per second)
	// 1: duration in seconds
	output, err := utils.ExecuteWithSudo(10*time.Second, "profile-bpfcc", "-F", "99", "1")
	if err != nil {
		// Even with error, we might have partial output
		if len(output) == 0 {
			return nil, err
		}
	}

	return parseProfileOutput(string(output)), nil
}

func parseProfileOutput(output string) []models.CPUProfile {
	// profile-bpfcc output format:
	// Stack trace lines followed by process name and count
	// Example:
	//     finish_task_switch
	//     schedule
	//     do_wait
	//     sys_wait4
	//     bash
	//         5
	//
	// Empty line separates different stack traces

	var profiles []models.CPUProfile
	lines := strings.Split(output, "\n")

	var currentStack []string
	var processName string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Skip empty lines and header
		if line == "" || strings.HasPrefix(line, "Sampling") {
			// If we have accumulated a stack, save it
			if len(currentStack) > 0 && processName != "" {
				// The last element before count is the process name
				stackTrace := strings.Join(currentStack, "\n")

				// Look for the count on the next line
				if i+1 < len(lines) {
					countLine := strings.TrimSpace(lines[i+1])
					if count, err := strconv.Atoi(countLine); err == nil {
						profile := models.CPUProfile{
							ProcessName: processName,
							StackTrace:  stackTrace,
							SampleCount: count,
						}
						profiles = append(profiles, profile)
						i++ // Skip the count line
					}
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
				profiles = append(profiles, profile)

				// Reset for next stack
				currentStack = []string{}
				processName = ""
			}
			continue
		}

		// This is a stack frame, add to current stack
		currentStack = append(currentStack, line)
	}

	return profiles
}
