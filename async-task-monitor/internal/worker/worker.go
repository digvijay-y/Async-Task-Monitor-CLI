package worker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"async-task-monitor/internal/model"
)

func RunTask(ctx context.Context, task *model.Task, updates chan<- model.Task) {
	task.Status = model.Running
	task.StartedAt = time.Now()
	task.Error = ""
	updates <- *task

	for i := 0; i <= 100; i += 10 {
		select {
		case <-ctx.Done():
			task.Status = model.Cancelled
			task.EndedAt = time.Now()
			task.Error = "cancelled by user"
			updates <- *task
			return
		case <-time.After(350 * time.Millisecond):
		}

		task.Progress = i
		task.Logs = append(task.Logs,
			fmt.Sprintf("progress reached %d%%", i))

		updates <- *task
	}

	if rand.Intn(10) < 2 {
		task.Status = model.Failed
		task.Error = "simulated worker failure"
	} else {
		task.Status = model.Success
	}

	task.EndedAt = time.Now()

	updates <- *task
}
