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
	actions := make(chan manager.Action, 16)

	p := tea.NewProgram(ui.InitialModel(actions))

	go func() {
		tasks := []model.Task{
			{ID: "1", Name: "Build API Service", Status: model.Queued},
			{ID: "2", Name: "Run Integration Tests", Status: model.Queued},
			{ID: "3", Name: "Deploy Production", Status: model.Queued},
			{ID: "4", Name: "Generate Reports", Status: model.Queued},
			{ID: "5", Name: "Sync S3 Backups", Status: model.Queued},
		}

		manager.RunTasks(tasks, updates, actions)

		for update := range updates {
			p.Send(ui.TaskUpdateMsg(update))
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
	}
}
