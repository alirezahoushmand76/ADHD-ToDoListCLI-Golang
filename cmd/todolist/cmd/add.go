package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/models"
	"github.com/user/todolist/internal/ui"
)

var (
	addTitle       string
	addDescription string
	addPriority    string
	addCategory    string
	addDueDate     string
	addReminder    string

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Long:  `Add a new task to your to-do list with title, description, priority, category, due date, and reminder.`,
		RunE:  runAddCmd,
		Example: `  todolist add "Complete project report" --priority high --category work --due tomorrow
  todolist add "Read book" --priority medium
  todolist add --title "Call doctor" --due "next week" --priority high`,
	}
)

func init() {
	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "Task title")
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Task description")
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "medium", "Task priority (low, medium, high)")
	addCmd.Flags().StringVarP(&addCategory, "category", "c", "inbox", "Task category")
	addCmd.Flags().StringVar(&addDueDate, "due", "", "Due date (YYYY-MM-DD, today, tomorrow, next week)")
	addCmd.Flags().StringVar(&addReminder, "reminder", "", "Reminder time (YYYY-MM-DD, today, tomorrow, next week)")
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// Parse due date if provided
	var dueDate time.Time
	if addDueDate != "" {
		var err error
		dueDate, err = parseDateTime(addDueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format: %w", err)
		}
	}

	// Parse reminder time if provided
	var reminderAt time.Time
	if addReminder != "" {
		var err error
		reminderAt, err = parseDateTime(addReminder)
		if err != nil {
			return fmt.Errorf("invalid reminder time format: %w", err)
		}
	}

	// Validate priority
	priority := models.Priority(strings.ToLower(addPriority))
	if priority != models.PriorityLow && priority != models.PriorityMedium && priority != models.PriorityHigh {
		return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", addPriority)
	}

	// Validate or create category
	category := models.Category(strings.ToLower(addCategory))

	var task *models.Task
	var err error

	if todoClient != nil {
		// Use the client to add the task
		task, err = todoClient.AddTask(
			title,
			addDescription,
			priority,
			category,
			dueDate,
			reminderAt,
		)
	} else {
		// Use the app directly (legacy mode)
		err = todoApp.AddTask(
			title,
			addDescription,
			priority,
			category,
			dueDate,
			reminderAt,
		)
		if err == nil {
			// Get the task that was just added (this is a bit of a hack)
			tasks, _ := todoApp.GetAllTasks()
			if len(tasks) > 0 {
				task = tasks[len(tasks)-1]
			}
		}
	}

	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}

	ui.PrintSuccess("Task added successfully!")
	if task != nil {
		fmt.Println(task.String())
	}
	return nil
}

// parseDateTime parses a date or date-time string
func parseDateTime(input string) (time.Time, error) {
	// Handle empty input
	if input == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Try parsing as date-time
	t, err := time.Parse("2006-01-02 15:04", input)
	if err == nil {
		return t, nil
	}

	// Try parsing as date
	t, err = time.Parse("2006-01-02", input)
	if err == nil {
		// Set time to end of day
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local), nil
	}

	// Try parsing as relative time
	return parseRelativeTime(input)
}

// parseRelativeTime parses relative time expressions like "today", "tomorrow", "next week"
func parseRelativeTime(input string) (time.Time, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	switch strings.ToLower(input) {
	case "today":
		return today, nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	case "next week":
		return today.AddDate(0, 0, 7), nil
	case "next month":
		return today.AddDate(0, 1, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown date format: %s (use YYYY-MM-DD or 'today', 'tomorrow', 'next week')", input)
	}
}
