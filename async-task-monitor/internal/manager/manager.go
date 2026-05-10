package manager

import (
	"context"
	"sync"
	"time"

	"async-task-monitor/internal/model"
	"async-task-monitor/internal/worker"
)

type ActionType string

const (
	ActionRetry  ActionType = "retry"
	ActionCancel ActionType = "cancel"
	ActionDelete ActionType = "delete"
)

type Action struct {
	Type   ActionType
	TaskID string
}

type Manager struct {
	workers int
	updates chan<- model.Task
	actions <-chan Action

	mu     sync.Mutex
	tasks  map[string]*model.Task
	active map[string]context.CancelFunc
	queue  []string
	closed bool
}

func New(workers int, updates chan<- model.Task, actions <-chan Action) *Manager {
	if workers < 1 {
		workers = 1
	}

	return &Manager{
		workers: workers,
		updates: updates,
		actions: actions,
		tasks:   make(map[string]*model.Task),
		active:  make(map[string]context.CancelFunc),
	}
}

func (m *Manager) Run(tasks []model.Task) {
	for i := range tasks {
		task := tasks[i]
		if task.MaxAttempts == 0 {
			task.MaxAttempts = 2
		}
		task.Status = model.Queued
		m.mu.Lock()
		m.tasks[task.ID] = &task
		m.queue = append(m.queue, task.ID)
		m.mu.Unlock()
		m.updates <- task
	}

	for {
		m.startQueued()

		if m.isIdle() && len(m.queue) == 0 {
			m.mu.Lock()
			if !m.closed {
				m.closed = true
				close(m.updates)
			}
			m.mu.Unlock()
			return
		}

		action, ok := <-m.actions
		if !ok {
			m.waitForActive()
			m.mu.Lock()
			if !m.closed {
				m.closed = true
				close(m.updates)
			}
			m.mu.Unlock()
			return
		}

		m.handleAction(action)
	}
}

func (m *Manager) handleAction(action Action) {
	m.mu.Lock()
	task, exists := m.tasks[action.TaskID]
	if !exists {
		m.mu.Unlock()
		return
	}

	switch action.Type {
	case ActionCancel:
		if cancel, ok := m.active[action.TaskID]; ok {
			cancel()
		}
	case ActionDelete:
		if _, ok := m.active[action.TaskID]; ok {
			m.mu.Unlock()
			return
		}
		delete(m.tasks, action.TaskID)
		m.mu.Unlock()
		m.updates <- model.Task{ID: action.TaskID, Status: model.Cancelled, Progress: -1, Error: "deleted"}
		return
	case ActionRetry:
		if _, ok := m.active[action.TaskID]; ok {
			m.mu.Unlock()
			return
		}
		if task.Status != model.Failed && task.Status != model.Cancelled {
			m.mu.Unlock()
			return
		}
		task.Status = model.Retrying
		task.Progress = 0
		task.Error = ""
		task.EndedAt = time.Time{}
		task.StartedAt = time.Time{}
		task.Attempt++
		m.queue = append([]string{action.TaskID}, m.queue...)
		m.mu.Unlock()
		m.updates <- *task
		return
	}

	m.mu.Unlock()
}

func (m *Manager) startQueued() {
	for {
		m.mu.Lock()
		if len(m.active) >= m.workers || len(m.queue) == 0 {
			m.mu.Unlock()
			return
		}

		taskID := m.queue[0]
		m.queue = m.queue[1:]
		task := m.tasks[taskID]
		ctx, cancel := context.WithCancel(context.Background())
		m.active[taskID] = cancel
		m.mu.Unlock()

		go func(task *model.Task, ctx context.Context, taskID string) {
			defer func() {
				m.mu.Lock()
				delete(m.active, taskID)
				m.mu.Unlock()
			}()

			if task.Attempt > 0 {
				task.Status = model.Retrying
			} else {
				task.Status = model.Running
			}

			worker.RunTask(ctx, task, m.updates)
		}(task, ctx, taskID)
	}
}

func (m *Manager) isIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.active) == 0
}

func (m *Manager) waitForActive() {
	for {
		if m.isIdle() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func RunTasks(tasks []model.Task, updates chan<- model.Task, actions <-chan Action) {
	m := New(4, updates, actions)
	m.Run(tasks)
}
