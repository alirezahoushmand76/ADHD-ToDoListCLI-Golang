package models

import (
	"fmt"
	"time"
)

// Priority represents the importance level of a task
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Category represents the grouping of a task
type Category string

// Task represents a to-do item
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Category    Category  `json:"category"`
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	ReminderAt  time.Time `json:"reminder_at"`
}

// NewTask creates a new task with the given parameters
func NewTask(title, description string, priority Priority, category Category, dueDate, reminderAt time.Time) *Task {
	return &Task{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()), // Simple ID generation
		Title:       title,
		Description: description,
		Priority:    priority,
		Category:    category,
		DueDate:     dueDate,
		Completed:   false,
		CreatedAt:   time.Now(),
		ReminderAt:  reminderAt,
	}
}

// IsOverdue checks if the task is past its due date
func (t *Task) IsOverdue() bool {
	return !t.DueDate.IsZero() && time.Now().After(t.DueDate) && !t.Completed
}

// IsReminderDue checks if it's time to remind about this task
func (t *Task) IsReminderDue() bool {
	return !t.ReminderAt.IsZero() && time.Now().After(t.ReminderAt) && !t.Completed
}

// MarkComplete marks the task as completed
func (t *Task) MarkComplete() {
	t.Completed = true
}

// String returns a string representation of the task
func (t *Task) String() string {
	status := "[ ]"
	if t.Completed {
		status = "[✓]"
	}

	dueStr := "No due date"
	if !t.DueDate.IsZero() {
		dueStr = t.DueDate.Format("2006-01-02 15:04")
	}

	return fmt.Sprintf("%s %s (Priority: %s, Category: %s, Due: %s)",
		status, t.Title, t.Priority, t.Category, dueStr)
}
