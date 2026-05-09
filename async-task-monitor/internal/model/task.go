package model

import "time"

type Status string

const (
	Pending   Status = "PENDING"
	Running   Status = "RUNNING"
	Success   Status = "SUCCESS"
	Failed    Status = "FAILED"
	Cancelled Status = "CANCELLED"
)

type Task struct {
	ID        string
	Name      string
	Status    Status
	Progress  int
	Logs      []string
	StartedAt time.Time
	EndedAt   time.Time
}
