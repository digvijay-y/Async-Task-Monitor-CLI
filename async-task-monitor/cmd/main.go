package main

import (
	"fmt"

	"async-task-monitor/internal/manager"
	"async-task-monitor/internal/model"
	"async-task-monitor/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	updates := make(chan model.Task, 16)

	p := tea.NewProgram(ui.InitialModel())

	go func() {
		tasks := []model.Task{
			{ID: "1", Name: "Build API", Status: model.Pending},
			{ID: "2", Name: "Run Tests", Status: model.Pending},
			{ID: "3", Name: "Deploy", Status: model.Pending},
		}

		manager.RunTasks(tasks, updates)

		for update := range updates {
			p.Send(ui.TaskUpdateMsg(update))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
	}
}