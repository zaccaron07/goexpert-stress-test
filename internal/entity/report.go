package entity

import "time"

type Report struct {
	URL           string        `json:"url"`
	TotalTime     time.Duration `json:"total_time"`
	TotalRequests int           `json:"total_requests"`
	StatusCodes   map[int]int   `json:"status_codes"`
}
