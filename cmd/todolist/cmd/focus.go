package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/ui"
)

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Enter focus mode",
	Long:  `Enter focus mode to get a suggestion for the next task to work on based on priority and urgency.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		task, err := todoApp.FocusMode()
		if err != nil {
			return fmt.Errorf("failed to enter focus mode: %w", err)
		}

		if task == nil {
			ui.PrintInfo("No tasks to focus on. Add some tasks first!")
			return nil
		}

		fmt.Println()
		ui.PrintInfo("🎯 FOCUS MODE")
		ui.PrintInfo("Here's the task you should focus on next:")
		fmt.Println()
		ui.PrintTask(task)
		fmt.Println()
		ui.PrintInfo("To start a Pomodoro timer for this task, run:")
		ui.PrintInfo("  todolist pomodoro %s", task.ID)
		fmt.Println()

		return nil
	},
}
