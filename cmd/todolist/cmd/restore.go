package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/ui"
)

var (
	restoreForce bool

	restoreCmd = &cobra.Command{
		Use:   "restore [backup_file_or_index]",
		Short: "Restore tasks from a backup",
		Long:  `Restore tasks from a backup file or by backup index.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			backupArg := args[0]

			// Check if the argument is an index
			index, err := strconv.Atoi(backupArg)
			if err == nil {
				// It's an index, get the backup file
				backups, err := todoApp.ListBackups()
				if err != nil {
					return fmt.Errorf("failed to list backups: %w", err)
				}

				if len(backups) == 0 {
					return fmt.Errorf("no backups found")
				}

				if index < 1 || index > len(backups) {
					return fmt.Errorf("invalid backup index: %d (must be between 1 and %d)", index, len(backups))
				}

				backupArg = backups[index-1]
			}

			// Confirm restore unless --force flag is used
			if !restoreForce {
				ui.PrintWarning("Are you sure you want to restore from backup: %s? This will replace all current tasks. (y/N): ", backupArg)
				var confirm string
				fmt.Scanln(&confirm)

				if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
					ui.PrintInfo("Restore cancelled")
					return nil
				}
			}

			// Restore from backup
			err = todoApp.RestoreTasks(backupArg)
			if err != nil {
				return fmt.Errorf("failed to restore tasks: %w", err)
			}

			ui.PrintSuccess("Tasks restored from: %s", backupArg)
			return nil
		},
		Example: `  todolist restore 1
  todolist restore /path/to/backup/file.json
  todolist restore 2 --force`,
	}
)

func init() {
	restoreCmd.Flags().BoolVarP(&restoreForce, "force", "f", false, "Restore without confirmation")
}
