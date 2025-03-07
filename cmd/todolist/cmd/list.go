package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/models"
	"github.com/user/todolist/internal/ui"
)

var (
	listCategory string
	listPriority string
	listAll      bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long:  `List tasks with optional filtering by category or priority.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var tasks []*models.Task
			var err error
			var title string

			// Get tasks based on filters
			if listCategory != "" {
				tasks, err = todoApp.GetTasksByCategory(models.Category(listCategory))
				if err != nil {
					return fmt.Errorf("failed to get tasks by category: %w", err)
				}
				title = fmt.Sprintf("Tasks in category: %s", listCategory)
			} else if listPriority != "" {
				// Validate priority
				var priority models.Priority
				switch strings.ToLower(listPriority) {
				case "low":
					priority = models.PriorityLow
				case "medium":
					priority = models.PriorityMedium
				case "high":
					priority = models.PriorityHigh
				default:
					return fmt.Errorf("invalid priority: %s (must be low, medium, or high)", listPriority)
				}

				tasks, err = todoApp.GetTasksByPriority(priority)
				if err != nil {
					return fmt.Errorf("failed to get tasks by priority: %w", err)
				}
				title = fmt.Sprintf("Tasks with priority: %s", priority)
			} else {
				tasks, err = todoApp.GetAllTasks()
				if err != nil {
					return fmt.Errorf("failed to get tasks: %w", err)
				}
				title = "All Tasks"
			}

			// Filter out completed tasks unless --all is specified
			if !listAll {
				var incompleteTasks []*models.Task
				for _, task := range tasks {
					if !task.Completed {
						incompleteTasks = append(incompleteTasks, task)
					}
				}
				tasks = incompleteTasks
				title += " (incomplete)"
			}

			// Print tasks
			ui.PrintTaskList(tasks, title)
			return nil
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter tasks by category")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter tasks by priority (low, medium, high)")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all tasks, including completed ones")
}
