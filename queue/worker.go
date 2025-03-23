package queue

import "time"

type Job struct {
	ID        string    `json:"id"`
	TimeStamp time.Time `json:"timestamp"`
}
