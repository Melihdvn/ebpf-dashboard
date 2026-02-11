package collector

import (
	"ebpf-dashboard/models"
	"ebpf-dashboard/utils"
	"regexp"
	"strconv"
	"time"
)

// CollectDiskLatency runs biolatency and returns parsed latency distribution
func CollectDiskLatency() ([]models.DiskLatency, error) {
	// Run biolatency with 1 second interval, 1 count
	output, err := utils.ExecuteWithSudo(5*time.Second, "biolatency", "1", "1")
	if err != nil {
		// Even with error, we might have partial output
		if len(output) == 0 {
			return nil, err
		}
	}

	return parseBiolatencyOutput(string(output)), nil
}

func parseBiolatencyOutput(output string) []models.DiskLatency {
	// biolatency output format (histogram):
	// usecs               : count     distribution
	//     0 -> 1          : 0        |                                        |
	//     2 -> 3          : 5        |*****                                   |
	//     4 -> 7          : 10       |**********                              |

	re := regexp.MustCompile(`(\d+)\s*->\s*(\d+)\s*:\s*(\d+)`)
	matches := re.FindAllStringSubmatch(output, -1)

	var latencies []models.DiskLatency
	for _, match := range matches {
		if len(match) == 4 {
			rangeMin, _ := strconv.Atoi(match[1])
			rangeMax, _ := strconv.Atoi(match[2])
			count, _ := strconv.Atoi(match[3])

			latency := models.DiskLatency{
				RangeMin: rangeMin,
				RangeMax: rangeMax,
				Count:    count,
			}
			latencies = append(latencies, latency)
		}
	}

	return latencies
}
