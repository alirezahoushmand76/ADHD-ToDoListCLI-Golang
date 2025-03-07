package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/todolist/internal/models"
	"github.com/user/todolist/internal/storage"
	"github.com/user/todolist/internal/ui"
	"github.com/user/todolist/internal/utils"
)

// App represents the todo list application
type App struct {
	Storage storage.Storage
	Config  *Config
}

// Config represents the application configuration
type Config struct {
	DataDir           string
	StorageFile       string
	BackupDir         string
	DefaultCategories []models.Category
}

// DefaultConfig returns the default application configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	dataDir := filepath.Join(homeDir, ".todolist")

	return &Config{
		DataDir:     dataDir,
		StorageFile: filepath.Join(dataDir, "tasks.json"),
		BackupDir:   filepath.Join(dataDir, "backups"),
		DefaultCategories: []models.Category{
			models.Category("work"),
			models.Category("personal"),
			models.Category("urgent"),
			models.Category("health"),
			models.Category("learning"),
		},
	}
}

// NewApp creates a new application instance
func NewApp(config *Config) (*App, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create storage
	store := storage.NewJSONStorage(config.StorageFile)
	if err := store.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return &App{
		Storage: store,
		Config:  config,
	}, nil
}

// AddTask adds a new task
func (a *App) AddTask(title, description string, priority models.Priority, category models.Category, dueDate, reminderAt time.Time) error {
	task := models.NewTask(title, description, priority, category, dueDate, reminderAt)
	return a.Storage.AddTask(task)
}

// GetTask retrieves a task by ID
func (a *App) GetTask(id string) (*models.Task, error) {
	return a.Storage.GetTask(id)
}

// GetAllTasks retrieves all tasks
func (a *App) GetAllTasks() ([]*models.Task, error) {
	return a.Storage.GetAllTasks()
}

// GetTasksByCategory retrieves tasks by category
func (a *App) GetTasksByCategory(category models.Category) ([]*models.Task, error) {
	return a.Storage.GetTasksByCategory(category)
}

// GetTasksByPriority retrieves tasks by priority
func (a *App) GetTasksByPriority(priority models.Priority) ([]*models.Task, error) {
	return a.Storage.GetTasksByPriority(priority)
}

// UpdateTask updates an existing task
func (a *App) UpdateTask(task *models.Task) error {
	return a.Storage.UpdateTask(task)
}

// DeleteTask deletes a task by ID
func (a *App) DeleteTask(id string) error {
	return a.Storage.DeleteTask(id)
}

// CompleteTask marks a task as completed
func (a *App) CompleteTask(id string) error {
	task, err := a.Storage.GetTask(id)
	if err != nil {
		return err
	}

	task.MarkComplete()
	return a.Storage.UpdateTask(task)
}

// BackupTasks creates a backup of the tasks
func (a *App) BackupTasks() (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(a.Config.BackupDir, fmt.Sprintf("tasks-backup-%s.json", timestamp))

	if err := a.Storage.Backup(backupFile); err != nil {
		return "", err
	}

	return backupFile, nil
}

// RestoreTasks restores tasks from a backup
func (a *App) RestoreTasks(backupFile string) error {
	return a.Storage.Restore(backupFile)
}

// ListBackups lists available backups
func (a *App) ListBackups() ([]string, error) {
	files, err := os.ReadDir(a.Config.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "tasks-backup-") && strings.HasSuffix(file.Name(), ".json") {
			backups = append(backups, filepath.Join(a.Config.BackupDir, file.Name()))
		}
	}

	return backups, nil
}

// BrainDump enters brain dump mode for quickly adding multiple tasks
func (a *App) BrainDump() error {
	scanner := bufio.NewScanner(os.Stdin)
	count := 0

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		title := strings.TrimSpace(scanner.Text())
		if title == "" {
			break
		}

		// Handle special commands
		if title == "q" || title == "quit" || title == "exit" {
			ui.PrintInfo("Exiting brain dump mode")
			break
		}

		if title == "help" || title == "--help" || title == "-h" {
			ui.PrintInfo("Brain Dump Mode Help:")
			ui.PrintInfo("- Enter task titles one per line")
			ui.PrintInfo("- Leave a line empty to finish")
			ui.PrintInfo("- Type 'q', 'quit', or 'exit' to exit")
			ui.PrintInfo("- Press Ctrl+C to exit at any time")
			continue
		}

		task := models.NewTask(title, "", models.PriorityMedium, models.Category("inbox"), time.Time{}, time.Time{})
		if err := a.Storage.AddTask(task); err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		count++
		ui.PrintSuccess("Added: %s", title)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	ui.PrintSuccess("Added %d tasks to your inbox.", count)
	return nil
}

// FocusMode enters focus mode to suggest the next task to work on
func (a *App) FocusMode() (*models.Task, error) {
	tasks, err := a.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	nextTask := utils.GetNextFocusTask(tasks)
	return nextTask, nil
}

// StartPomodoro starts a Pomodoro timer for a task
func (a *App) StartPomodoro(taskID string, customDuration time.Duration) error {
	task, err := a.GetTask(taskID)
	if err != nil {
		return err
	}

	config := utils.DefaultPomodoroConfig()
	config.TaskName = task.Title

	if customDuration > 0 {
		config.WorkDuration = customDuration
	}

	session := utils.NewPomodoroSession(config)
	session.Start()

	return nil
}
