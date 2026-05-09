package ui

import (
	"fmt"
	"strings"

	"async-task-monitor/internal/model"

	tea "charm.land/bubbletea/v2"
)

type TaskUpdateMsg model.Task

type Model struct {
	Tasks map[string]model.Task
}

func InitialModel() Model {
	return Model{
		Tasks: make(map[string]model.Task),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TaskUpdateMsg:
		task := model.Task(msg)
		m.Tasks[task.ID] = task
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("Async Task Monitor\n\n")

	for _, t := range m.Tasks {
		b.WriteString(fmt.Sprintf(
			"%s | %-15s | %-10s | %d%%\n",
			t.ID,
			t.Name,
			t.Status,
			t.Progress,
		))
	}

	return b.String()
}