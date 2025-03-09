package protocol

import (
	"encoding/json"
	"time"

	"github.com/user/todolist/internal/models"
)

// Operation types
const (
	// Task operations
	OpAddTask            = "ADD_TASK"
	OpGetTask            = "GET_TASK"
	OpGetAllTasks        = "GET_ALL_TASKS"
	OpGetTasksByCategory = "GET_TASKS_BY_CATEGORY"
	OpGetTasksByPriority = "GET_TASKS_BY_PRIORITY"
	OpUpdateTask         = "UPDATE_TASK"
	OpDeleteTask         = "DELETE_TASK"
	OpCompleteTask       = "COMPLETE_TASK"

	// Data operations
	OpBackup      = "BACKUP"
	OpRestore     = "RESTORE"
	OpListBackups = "LIST_BACKUPS"

	// Other operations
	OpBrainDump     = "BRAIN_DUMP"
	OpFocusMode     = "FOCUS_MODE"
	OpStartPomodoro = "START_POMODORO"
)

// Request represents a client request to the server
type Request struct {
	Operation string          `json:"operation"`
	Payload   json.RawMessage `json:"payload"`
}

// Response represents a server response to the client
type Response struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// AddTaskRequest represents the payload for adding a task
type AddTaskRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Priority    models.Priority `json:"priority"`
	Category    models.Category `json:"category"`
	DueDate     time.Time       `json:"due_date"`
	ReminderAt  time.Time       `json:"reminder_at"`
}

// TaskResponse represents a task in a response
type TaskResponse struct {
	Task *models.Task `json:"task"`
}

// TasksResponse represents multiple tasks in a response
type TasksResponse struct {
	Tasks []*models.Task `json:"tasks"`
}

// IDRequest represents a request with just an ID
type IDRequest struct {
	ID string `json:"id"`
}

// CategoryRequest represents a request with a category
type CategoryRequest struct {
	Category models.Category `json:"category"`
}

// PriorityRequest represents a request with a priority
type PriorityRequest struct {
	Priority models.Priority `json:"priority"`
}

// BackupRequest represents a request to backup data
type BackupRequest struct {
	Filename string `json:"filename,omitempty"`
}

// BackupResponse represents the response to a backup request
type BackupResponse struct {
	Filename string `json:"filename"`
}

// RestoreRequest represents a request to restore data
type RestoreRequest struct {
	Filename string `json:"filename"`
}

// ListBackupsResponse represents the response to a list backups request
type ListBackupsResponse struct {
	Backups []string `json:"backups"`
}

// BrainDumpResponse represents the response to a brain dump request
type BrainDumpResponse struct {
	Success bool `json:"success"`
}

// PomodoroRequest represents a request to start a pomodoro timer
type PomodoroRequest struct {
	TaskID         string        `json:"task_id"`
	CustomDuration time.Duration `json:"custom_duration,omitempty"`
}

// StringResponse represents a simple string response
type StringResponse struct {
	Message string `json:"message"`
}
