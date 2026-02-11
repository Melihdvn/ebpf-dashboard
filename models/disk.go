package models

import "time"

type DiskLatency struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	RangeMin  int       `json:"range_min"`
	RangeMax  int       `json:"range_max"`
	Count     int       `json:"count"`
}
