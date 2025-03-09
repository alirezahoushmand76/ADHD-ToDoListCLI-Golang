package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/models"
)

var (
	listCategory string
	listPriority string
	listAll      bool
	listVerbose  bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long:  `List tasks with optional filtering by category or priority.`,
		RunE:  runListCmd,
	}
)

func init() {
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter tasks by category")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter tasks by priority (low, medium, high)")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all tasks, including completed ones")
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show detailed task information")
}

func runListCmd(cmd *cobra.Command, args []string) error {
	var tasks []*models.Task
	var err error

	if todoClient != nil {
		// Use the client to get tasks
		if listCategory != "" {
			tasks, err = todoClient.GetTasksByCategory(models.Category(listCategory))
		} else if listPriority != "" {
			priority := models.Priority(strings.ToLower(listPriority))
			if priority != models.PriorityLow && priority != models.PriorityMedium && priority != models.PriorityHigh {
				return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", listPriority)
			}
			tasks, err = todoClient.GetTasksByPriority(priority)
		} else {
			tasks, err = todoClient.GetAllTasks()
		}
	} else {
		// Use the app directly (legacy mode)
		if listCategory != "" {
			tasks, err = todoApp.GetTasksByCategory(models.Category(listCategory))
		} else if listPriority != "" {
			priority := models.Priority(strings.ToLower(listPriority))
			if priority != models.PriorityLow && priority != models.PriorityMedium && priority != models.PriorityHigh {
				return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", listPriority)
			}
			tasks, err = todoApp.GetTasksByPriority(priority)
		} else {
			tasks, err = todoApp.GetAllTasks()
		}
	}

	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}

	// Filter completed tasks if requested
	if !listAll {
		var filteredTasks []*models.Task
		for _, task := range tasks {
			if !task.Completed {
				filteredTasks = append(filteredTasks, task)
			}
		}
		tasks = filteredTasks
	}

	// Display tasks
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Printf("Found %d tasks:\n\n", len(tasks))
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task.String())
		if listVerbose {
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			fmt.Printf("   ID: %s\n", task.ID)
			fmt.Printf("   Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
			if !task.ReminderAt.IsZero() {
				fmt.Printf("   Reminder: %s\n", task.ReminderAt.Format("2006-01-02 15:04"))
			}
			fmt.Println()
		}
	}

	return nil
}
