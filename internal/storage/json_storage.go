package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/user/todolist/internal/models"
)

// JSONStorage implements the Storage interface using JSON files
type JSONStorage struct {
	filePath string
	tasks    map[string]*models.Task
	mu       sync.RWMutex
}

// NewJSONStorage creates a new JSON storage with the given file path
func NewJSONStorage(filePath string) *JSONStorage {
	return &JSONStorage{
		filePath: filePath,
		tasks:    make(map[string]*models.Task),
	}
}

// Initialize loads tasks from the JSON file or creates a new file if it doesn't exist
func (s *JSONStorage) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Create empty file
		s.tasks = make(map[string]*models.Task)
		return s.saveToFile()
	}

	// Read file
	file, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode JSON
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		s.tasks = make(map[string]*models.Task)
		return nil
	}

	var tasks []*models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Convert to map
	s.tasks = make(map[string]*models.Task)
	for _, task := range tasks {
		s.tasks[task.ID] = task
	}

	return nil
}

// saveToFile saves tasks to the JSON file
func (s *JSONStorage) saveToFile() error {
	// Convert map to slice
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	// Encode JSON
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// AddTask adds a new task
func (s *JSONStorage) AddTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return s.saveToFile()
}

// GetTask retrieves a task by ID
func (s *JSONStorage) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound{ID: id}
	}
	return task, nil
}

// GetAllTasks retrieves all tasks
func (s *JSONStorage) GetAllTasks() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTasksByCategory retrieves tasks by category
func (s *JSONStorage) GetTasksByCategory(category models.Category) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*models.Task
	for _, task := range s.tasks {
		if task.Category == category {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// GetTasksByPriority retrieves tasks by priority
func (s *JSONStorage) GetTasksByPriority(priority models.Priority) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*models.Task
	for _, task := range s.tasks {
		if task.Priority == priority {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// UpdateTask updates an existing task
func (s *JSONStorage) UpdateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return ErrTaskNotFound{ID: task.ID}
	}

	s.tasks[task.ID] = task
	return s.saveToFile()
}

// DeleteTask deletes a task by ID
func (s *JSONStorage) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound{ID: id}
	}

	delete(s.tasks, id)
	return s.saveToFile()
}

// Backup creates a backup of the tasks
func (s *JSONStorage) Backup(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	// Encode JSON
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Write to backup file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// Restore restores tasks from a backup
func (s *JSONStorage) Restore(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read backup file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var tasks []*models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Convert to map
	s.tasks = make(map[string]*models.Task)
	for _, task := range tasks {
		s.tasks[task.ID] = task
	}

	// Save to file
	return s.saveToFile()
}
