package task

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

const maxTasks = 10

// Manager allows us to operate on tasks.
type Manager struct {
	mu    sync.RWMutex
	tasks []Task
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{
		tasks: make([]Task, 0),
	}
}

// Add adds a new Task.
func (m *Manager) Add(content string) (*Task, error) {
	if content == "" {
		return nil, errors.New("task content cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, task := range m.tasks {
		if task.DoneAt.IsZero() && task.Content == content {
			return nil, errors.New("task already exists")
		}
	}
	task := Task{
		ID:        uuid.NewV4(),
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	m.tasks = append(m.tasks, task)
	if len(m.tasks) > maxTasks {
		m.tasks = m.tasks[1:]
	}

	return &task, nil
}

// MarkAsDone marks a new Task as done.
func (m *Manager) MarkAsDone(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var found bool
	for i, task := range m.tasks {
		if task.ID == id {
			m.tasks[i].DoneAt = time.Now().UTC()
			found = true
			break
		}
	}
	if !found {
		return errors.New("task not found")
	}

	return nil
}

// List returns a list of non-done tasks.
func (m *Manager) List() []Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	filtered := lo.Filter(m.tasks, func(task Task, _ int) bool {
		return task.DoneAt.IsZero()
	})

	slices.SortFunc(filtered, func(a, b Task) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return 1
		}

		return -1
	})

	return filtered
}

// Seed adds number num of random tasks.
func (m *Manager) Seed(num uint8) error {
	for range num {
		if _, err := m.Add(gofakeit.SentenceSimple()); err != nil {
			return err
		}
	}

	return nil
}
