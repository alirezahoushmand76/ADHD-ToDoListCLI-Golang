package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/storage"
	"github.com/user/todolist/internal/ui"
)

var (
	deleteForce bool

	deleteCmd = &cobra.Command{
		Use:   "delete [task_id]",
		Short: "Delete a task",
		Long:  `Delete a task by its ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

			// Validate task ID format
			if strings.Contains(taskID, "[") || strings.Contains(taskID, "]") {
				return fmt.Errorf("invalid task ID format: %s (do not include square brackets)", taskID)
			}

			// Get the task first to check if it exists and to show the title in the success message
			task, err := todoApp.GetTask(taskID)
			if err != nil {
				// Provide a more helpful error message for task not found
				if _, ok := err.(storage.ErrTaskNotFound); ok {
					return fmt.Errorf("task not found with ID: %s (use 'todolist list' to see all tasks)", taskID)
				}
				return fmt.Errorf("failed to get task: %w", err)
			}

			// Confirm deletion unless --force flag is used
			if !deleteForce {
				ui.PrintWarning("Are you sure you want to delete task: %s? (y/N): ", task.Title)
				var confirm string
				fmt.Scanln(&confirm)

				if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
					ui.PrintInfo("Task deletion cancelled")
					return nil
				}
			}

			// Delete the task
			err = todoApp.DeleteTask(taskID)
			if err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			ui.PrintSuccess("Task deleted: %s", task.Title)
			return nil
		},
		Example: `  todolist delete 1741359296120413000
  todolist delete 1741359296120413000 --force`,
	}
)

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Delete without confirmation")
}
