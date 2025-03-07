package storage

import (
	"github.com/user/todolist/internal/models"
)

// Storage defines the interface for task persistence
type Storage interface {
	// Task operations
	AddTask(task *models.Task) error
	GetTask(id string) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	GetTasksByCategory(category models.Category) ([]*models.Task, error)
	GetTasksByPriority(priority models.Priority) ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error

	// Data operations
	Backup(filename string) error
	Restore(filename string) error

	// Initialize storage
	Initialize() error
}

// ErrTaskNotFound is returned when a task with the specified ID is not found
type ErrTaskNotFound struct {
	ID string
}

func (e ErrTaskNotFound) Error() string {
	return "task not found: " + e.ID
}
