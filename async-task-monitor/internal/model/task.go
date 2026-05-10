package model

import "time"

type Status string

const (
	Queued    Status = "QUEUED"
	Pending   Status = "PENDING"
	Running   Status = "RUNNING"
	Retrying  Status = "RETRYING"
	Success   Status = "SUCCESS"
	Failed    Status = "FAILED"
	Cancelled Status = "CANCELLED"
)

type Task struct {
	ID          string
	Name        string
	Status      Status
	Progress    int
	Logs        []string
	StartedAt   time.Time
	EndedAt     time.Time
	Attempt     int
	MaxAttempts int
	Error       string
}
