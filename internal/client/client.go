package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/user/todolist/internal/models"
	"github.com/user/todolist/internal/protocol"
)

// Client represents a connection to the TodoList server
type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

// NewClient creates a new client connected to the specified address
func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

// Close closes the connection to the server
func (c *Client) Close() error {
	return c.conn.Close()
}

// sendRequest sends a request to the server and returns the response
func (c *Client) sendRequest(operation string, payload interface{}) (*protocol.Response, error) {
	// Marshal payload
	var payloadBytes []byte
	var err error
	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	// Create request
	request := protocol.Request{
		Operation: operation,
		Payload:   payloadBytes,
	}

	// Marshal request
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	requestBytes = append(requestBytes, '\n')

	// Send request
	_, err = c.conn.Write(requestBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	responseBytes, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal response
	var response protocol.Response
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// AddTask adds a new task
func (c *Client) AddTask(title, description string, priority models.Priority, category models.Category, dueDate, reminderAt time.Time) (*models.Task, error) {
	payload := protocol.AddTaskRequest{
		Title:       title,
		Description: description,
		Priority:    priority,
		Category:    category,
		DueDate:     dueDate,
		ReminderAt:  reminderAt,
	}

	response, err := c.sendRequest(protocol.OpAddTask, payload)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var taskResp protocol.TaskResponse
	if err := json.Unmarshal(response.Payload, &taskResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
	}

	return taskResp.Task, nil
}

// GetTask retrieves a task by ID
func (c *Client) GetTask(id string) (*models.Task, error) {
	payload := protocol.IDRequest{ID: id}

	response, err := c.sendRequest(protocol.OpGetTask, payload)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var taskResp protocol.TaskResponse
	if err := json.Unmarshal(response.Payload, &taskResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
	}

	return taskResp.Task, nil
}

// GetAllTasks retrieves all tasks
func (c *Client) GetAllTasks() ([]*models.Task, error) {
	response, err := c.sendRequest(protocol.OpGetAllTasks, nil)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var tasksResp protocol.TasksResponse
	if err := json.Unmarshal(response.Payload, &tasksResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks response: %w", err)
	}

	return tasksResp.Tasks, nil
}

// GetTasksByCategory retrieves tasks by category
func (c *Client) GetTasksByCategory(category models.Category) ([]*models.Task, error) {
	payload := protocol.CategoryRequest{Category: category}

	response, err := c.sendRequest(protocol.OpGetTasksByCategory, payload)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var tasksResp protocol.TasksResponse
	if err := json.Unmarshal(response.Payload, &tasksResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks response: %w", err)
	}

	return tasksResp.Tasks, nil
}

// GetTasksByPriority retrieves tasks by priority
func (c *Client) GetTasksByPriority(priority models.Priority) ([]*models.Task, error) {
	payload := protocol.PriorityRequest{Priority: priority}

	response, err := c.sendRequest(protocol.OpGetTasksByPriority, payload)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var tasksResp protocol.TasksResponse
	if err := json.Unmarshal(response.Payload, &tasksResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks response: %w", err)
	}

	return tasksResp.Tasks, nil
}

// UpdateTask updates a task
func (c *Client) UpdateTask(task *models.Task) error {
	payload := protocol.TaskResponse{Task: task}

	response, err := c.sendRequest(protocol.OpUpdateTask, payload)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// DeleteTask deletes a task by ID
func (c *Client) DeleteTask(id string) error {
	payload := protocol.IDRequest{ID: id}

	response, err := c.sendRequest(protocol.OpDeleteTask, payload)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// CompleteTask marks a task as completed
func (c *Client) CompleteTask(id string) error {
	payload := protocol.IDRequest{ID: id}

	response, err := c.sendRequest(protocol.OpCompleteTask, payload)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// BackupTasks creates a backup of all tasks
func (c *Client) BackupTasks() (string, error) {
	payload := protocol.BackupRequest{}

	response, err := c.sendRequest(protocol.OpBackup, payload)
	if err != nil {
		return "", err
	}

	if !response.Success {
		return "", fmt.Errorf("server error: %s", response.Error)
	}

	var backupResp protocol.BackupResponse
	if err := json.Unmarshal(response.Payload, &backupResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal backup response: %w", err)
	}

	return backupResp.Filename, nil
}

// RestoreTasks restores tasks from a backup
func (c *Client) RestoreTasks(filename string) error {
	payload := protocol.RestoreRequest{Filename: filename}

	response, err := c.sendRequest(protocol.OpRestore, payload)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// ListBackups lists all available backups
func (c *Client) ListBackups() ([]string, error) {
	response, err := c.sendRequest(protocol.OpListBackups, nil)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var backupsResp protocol.ListBackupsResponse
	if err := json.Unmarshal(response.Payload, &backupsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backups response: %w", err)
	}

	return backupsResp.Backups, nil
}

// BrainDump performs a brain dump
func (c *Client) BrainDump() error {
	response, err := c.sendRequest(protocol.OpBrainDump, nil)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// FocusMode enters focus mode
func (c *Client) FocusMode() (*models.Task, error) {
	response, err := c.sendRequest(protocol.OpFocusMode, nil)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("server error: %s", response.Error)
	}

	var taskResp protocol.TaskResponse
	if err := json.Unmarshal(response.Payload, &taskResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
	}

	return taskResp.Task, nil
}

// StartPomodoro starts a pomodoro timer for a task
func (c *Client) StartPomodoro(taskID string, customDuration time.Duration) error {
	payload := protocol.PomodoroRequest{
		TaskID:         taskID,
		CustomDuration: customDuration,
	}

	response, err := c.sendRequest(protocol.OpStartPomodoro, payload)
	if err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}
