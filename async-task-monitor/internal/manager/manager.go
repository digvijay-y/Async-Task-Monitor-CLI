package manager

import (
	"sync"

	"async-task-monitor/internal/model"
	"async-task-monitor/internal/worker"
)

func RunTasks(tasks []model.Task, updates chan<- model.Task) {
	var wg sync.WaitGroup

	for i := range tasks {
		wg.Add(1)

		go func(task *model.Task) {
			defer wg.Done()
			worker.RunTask(task, updates)
		}(&tasks[i])
	}

	wg.Wait()
	close(updates)
}
