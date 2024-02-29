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
func (m *Manager) Add(content string) error {
	if content == "" {
		return errors.New("task content cannot be empty")
	}

	m.mu.Lock()
	m.tasks = append(m.tasks, Task{
		ID:        uuid.NewV4(),
		Content:   content,
		CreatedAt: time.Now().UTC(),
	})
	if len(m.tasks) > maxTasks {
		m.tasks = m.tasks[1:]
	}
	m.mu.Unlock()

	return nil
}

// MarkAsDone marks a new Task as done.
func (m *Manager) MarkAsDone(id uuid.UUID) error {
	m.mu.Lock()
	var found bool
	for i, task := range m.tasks {
		if task.ID == id {
			m.tasks[i].DoneAt = time.Now().UTC()
			found = true
			break
		}
	}
	m.mu.Unlock()
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
		if err := m.Add(gofakeit.SentenceSimple()); err != nil {
			return err
		}
	}

	return nil
}
