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
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no title is provided, use the first argument
			if addTitle == "" && len(args) > 0 {
				addTitle = args[0]
			}

			// Title is required
			if addTitle == "" {
				return fmt.Errorf("title is required (use --title or provide it as an argument)")
			}

			// Parse priority with validation
			var priority models.Priority
			switch strings.ToLower(addPriority) {
			case "low":
				priority = models.PriorityLow
			case "medium":
				priority = models.PriorityMedium
			case "high":
				priority = models.PriorityHigh
			default:
				return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", addPriority)
			}

			// Parse category
			category := models.Category(addCategory)
			if category == "" {
				category = models.Category("inbox")
			}

			// Parse due date
			var dueDate time.Time
			if addDueDate != "" {
				var err error
				dueDate, err = parseDateTime(addDueDate)
				if err != nil {
					return fmt.Errorf("invalid due date: %w (use format YYYY-MM-DD or 'today', 'tomorrow', 'next week')", err)
				}
			}

			// Parse reminder
			var reminderAt time.Time
			if addReminder != "" {
				var err error
				reminderAt, err = parseDateTime(addReminder)
				if err != nil {
					return fmt.Errorf("invalid reminder: %w (use format YYYY-MM-DD or 'today', 'tomorrow', 'next week')", err)
				}
			}

			// Add task
			err := todoApp.AddTask(addTitle, addDescription, priority, category, dueDate, reminderAt)
			if err != nil {
				return fmt.Errorf("failed to add task: %w", err)
			}

			ui.PrintSuccess("Task added successfully: %s", addTitle)
			return nil
		},
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
