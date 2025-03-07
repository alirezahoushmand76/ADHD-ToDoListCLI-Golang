package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/storage"
	"github.com/user/todolist/internal/ui"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task_id]",
	Short: "Mark a task as completed",
	Long:  `Mark a task as completed by its ID.`,
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

		// Check if task is already completed
		if task.Completed {
			ui.PrintInfo("Task '%s' is already marked as completed", task.Title)
			return nil
		}

		// Mark as completed
		err = todoApp.CompleteTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		ui.PrintSuccess("Task completed: %s", task.Title)
		return nil
	},
	Example: `  todolist complete 1741359296120413000`,
}
