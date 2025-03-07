package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/storage"
	"github.com/user/todolist/internal/ui"
)

var (
	pomodoroDuration int

	pomodoroCmd = &cobra.Command{
		Use:   "pomodoro [task_id]",
		Short: "Start a Pomodoro timer for a task",
		Long:  `Start a Pomodoro timer for a task to help you focus on it.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			// Validate task ID format
			if strings.Contains(taskID, "[") || strings.Contains(taskID, "]") {
				return fmt.Errorf("invalid task ID format: %s (do not include square brackets)", taskID)
			}

			// Get the task first to check if it exists
			task, err := todoApp.GetTask(taskID)
			if err != nil {
				// Provide a more helpful error message for task not found
				if _, ok := err.(storage.ErrTaskNotFound); ok {
					return fmt.Errorf("task not found with ID: %s (use 'todolist list' to see all tasks)", taskID)
				}
				return fmt.Errorf("failed to get task: %w", err)
			}

			// Validate duration
			if pomodoroDuration < 0 {
				return fmt.Errorf("invalid duration: %d (must be a positive number)", pomodoroDuration)
			}

			// Inform user about how to exit
			ui.PrintInfo("Starting Pomodoro timer for task: %s", task.Title)
			ui.PrintInfo("Press Ctrl+C at any time to exit the timer")
			fmt.Println()

			// Start Pomodoro timer
			var duration time.Duration
			if pomodoroDuration > 0 {
				duration = time.Duration(pomodoroDuration) * time.Minute
			}

			return todoApp.StartPomodoro(taskID, duration)
		},
		Example: `  todolist pomodoro 1741359296120413000
  todolist pomodoro 1741359296120413000 --duration 30`,
	}
)

func init() {
	pomodoroCmd.Flags().IntVarP(&pomodoroDuration, "duration", "d", 0, "Custom work duration in minutes (default: 25)")
}
