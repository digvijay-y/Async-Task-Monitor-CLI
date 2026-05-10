package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"async-task-monitor/internal/manager"
	"async-task-monitor/internal/model"

	tea "charm.land/bubbletea/v2"
)

type TaskUpdateMsg model.Task

type panel int

const (
	panelOverview panel = iota
	panelLogs
)

type Model struct {
	Tasks    map[string]model.Task
	order    []string
	selected int
	panel    panel
	width    int
	height   int
	status   string
	ready    bool
	actions  chan<- manager.Action
	showHelp bool
	quitting bool
}

func InitialModel(actions chan<- manager.Action) Model {
	return Model{
		Tasks:   make(map[string]model.Task),
		order:   make([]string, 0),
		actions: actions,
		status:  "q quit  ↑↓ navigate  tab switch panel  enter logs  r retry  c cancel  d delete  h help",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case TaskUpdateMsg:
		task := model.Task(msg)
		if task.Progress < 0 {
			delete(m.Tasks, task.ID)
			m.rebuildOrder()
			m.clampSelection()
			return m, nil
		}
		m.Tasks[task.ID] = task
		m.rebuildOrder()
		m.clampSelection()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			m.moveSelection(-1)
		case "down", "j":
			m.moveSelection(1)
		case "tab":
			if m.panel == panelOverview {
				m.panel = panelLogs
			} else {
				m.panel = panelOverview
			}
		case "enter":
			m.panel = panelLogs
		case "h":
			m.showHelp = !m.showHelp
		case "r":
			if id, ok := m.selectedTaskID(); ok {
				return m, m.sendAction(manager.Action{Type: manager.ActionRetry, TaskID: id})
			}
		case "c":
			if id, ok := m.selectedTaskID(); ok {
				return m, m.sendAction(manager.Action{Type: manager.ActionCancel, TaskID: id})
			}
		case "d":
			if id, ok := m.selectedTaskID(); ok {
				if task, exists := m.Tasks[id]; exists && task.Status != model.Running && task.Status != model.Queued && task.Status != model.Pending && task.Status != model.Retrying {
					return m, m.sendAction(manager.Action{Type: manager.ActionDelete, TaskID: id})
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading Async Task Monitor...\n")
	}

	var b strings.Builder
	b.WriteString("Async Task Monitor CLI\n")
	b.WriteString(strings.Repeat("=", 80))
	b.WriteString("\n")
	b.WriteString(m.summaryLine())
	b.WriteString("\n\n")
	b.WriteString(m.taskTable())
	b.WriteString("\n\n")
	b.WriteString(m.logPanel())
	if m.showHelp {
		b.WriteString("\n\n")
		b.WriteString(box("Help", m.helpText()))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 80))
	b.WriteString("\n")
	b.WriteString(m.statusLine())

	v := tea.NewView(b.String())
	v.AltScreen = true
	v.WindowTitle = "Async Task Monitor CLI"
	return v
}

func (m Model) sendAction(action manager.Action) tea.Cmd {
	return func() tea.Msg {
		m.actions <- action
		return nil
	}
}

func (m *Model) rebuildOrder() {
	m.order = m.order[:0]
	for id := range m.Tasks {
		m.order = append(m.order, id)
	}
	sort.Slice(m.order, func(i, j int) bool {
		a := m.Tasks[m.order[i]]
		b := m.Tasks[m.order[j]]
		if a.StartedAt.Equal(b.StartedAt) {
			return a.ID < b.ID
		}
		if a.StartedAt.IsZero() {
			return false
		}
		if b.StartedAt.IsZero() {
			return true
		}
		return a.StartedAt.Before(b.StartedAt)
	})
	if len(m.order) == 0 {
		m.selected = 0
		return
	}
	if m.selected >= len(m.order) {
		m.selected = len(m.order) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

func (m *Model) clampSelection() {
	if len(m.order) == 0 {
		m.selected = 0
		return
	}
	if m.selected >= len(m.order) {
		m.selected = len(m.order) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

func (m *Model) moveSelection(delta int) {
	if len(m.order) == 0 {
		return
	}
	m.selected += delta
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= len(m.order) {
		m.selected = len(m.order) - 1
	}
}

func (m Model) selectedTaskID() (string, bool) {
	if len(m.order) == 0 || m.selected < 0 || m.selected >= len(m.order) {
		return "", false
	}
	return m.order[m.selected], true
}

func (m Model) summaryLine() string {
	var queued, running, success, failed, cancelled int
	for _, task := range m.Tasks {
		switch task.Status {
		case model.Queued, model.Pending:
			queued++
		case model.Running, model.Retrying:
			running++
		case model.Success:
			success++
		case model.Failed:
			failed++
		case model.Cancelled:
			cancelled++
		}
	}

	return fmt.Sprintf(
		"Tasks: %d    Running: %d    Success: %d    Failed: %d    Cancelled: %d    Queued: %d",
		len(m.Tasks), running, success, failed, cancelled, queued,
	)
}

func (m Model) taskTable() string {
	if len(m.order) == 0 {
		return box("Tasks", "No tasks yet. The queue will populate as work starts.")
	}

	var rows []string
	rows = append(rows, fmt.Sprintf("%-4s %-24s %-12s %-12s %-10s", "ID", "TASK NAME", "STATUS", "PROGRESS", "DURATION"))
	rows = append(rows, strings.Repeat("-", 70))
	for i, id := range m.order {
		task := m.Tasks[id]
		line := fmt.Sprintf(
			"%-4s %-24s %-12s %-12s %-10s",
			task.ID,
			truncate(task.Name, 24),
			string(task.Status),
			progressBar(task.Progress),
			formatDuration(task),
		)
		if i == m.selected {
			line = "> " + line
		} else {
			line = "  " + line
		}
		rows = append(rows, line)
	}

	return box("Tasks", strings.Join(rows, "\n"))
}

func (m Model) logPanel() string {
	if len(m.order) == 0 {
		return ""
	}

	id, ok := m.selectedTaskID()
	if !ok {
		return ""
	}
	task := m.Tasks[id]

	var lines []string
	lines = append(lines, fmt.Sprintf("Logs for task %s", task.ID))
	if len(task.Logs) == 0 {
		lines = append(lines, "No logs yet.")
		return box("Logs", strings.Join(lines, "\n"))
	}

	start := 0
	if len(task.Logs) > 6 {
		start = len(task.Logs) - 6
	}
	lines = append(lines, task.Logs[start:]...)
	return box("Logs", strings.Join(lines, "\n"))
}

func (m Model) helpText() string {
	return strings.Join([]string{
		"q or ctrl+c quit",
		"up/down navigate tasks",
		"tab switch panels",
		"enter open logs panel",
		"r retry failed task",
		"c cancel running task",
		"d delete completed task",
		"h toggle this help",
	}, "\n")
}

func (m Model) statusLine() string {
	if m.quitting {
		return "shutting down"
	}
	if len(m.order) == 0 {
		return m.status
	}
	id, _ := m.selectedTaskID()
	task := m.Tasks[id]
	return fmt.Sprintf("selected %s | %s | attempt %d", task.ID, task.Status, task.Attempt+1)
}

func progressBar(progress int) string {
	const width = 10
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	filled := (progress * width) / 100
	if filled > width {
		filled = width
	}
	return fmt.Sprintf("%s %3d%%", strings.Repeat("█", filled)+strings.Repeat("░", width-filled), progress)
}

func formatDuration(task model.Task) string {
	if task.StartedAt.IsZero() {
		return "--"
	}
	end := task.EndedAt
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(task.StartedAt).Truncate(time.Second).String()
}

func truncate(value string, width int) string {
	if len(value) <= width {
		return value
	}
	if width < 2 {
		return value[:width]
	}
	return value[:width-1] + "…"
}

func box(title, content string) string {
	lines := strings.Split(content, "\n")
	width := len(title) + 4
	for _, line := range lines {
		if len(line)+4 > width {
			width = len(line) + 4
		}
	}

	top := "+" + strings.Repeat("-", width-2) + "+"
	head := fmt.Sprintf("| %-*s |", width-4, title)
	sep := "+" + strings.Repeat("-", width-2) + "+"
	var b strings.Builder
	b.WriteString(top)
	b.WriteString("\n")
	b.WriteString(head)
	b.WriteString("\n")
	b.WriteString(sep)
	b.WriteString("\n")
	for _, line := range lines {
		b.WriteString(fmt.Sprintf("| %-*s |", width-4, line))
		b.WriteString("\n")
	}
	b.WriteString(top)
	return b.String()
}
