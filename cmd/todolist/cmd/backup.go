package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/ui"
)

var (
	listBackups bool

	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup tasks",
		Long:  `Create a backup of your tasks or list existing backups.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if listBackups {
				// List backups
				backups, err := todoApp.ListBackups()
				if err != nil {
					return fmt.Errorf("failed to list backups: %w", err)
				}

				if len(backups) == 0 {
					ui.PrintInfo("No backups found.")
					return nil
				}

				fmt.Println("Available backups:")
				for i, backup := range backups {
					fmt.Printf("%d. %s\n", i+1, backup)
				}
				return nil
			}

			// Create backup
			backupFile, err := todoApp.BackupTasks()
			if err != nil {
				return fmt.Errorf("failed to backup tasks: %w", err)
			}

			ui.PrintSuccess("Tasks backed up to: %s", backupFile)
			return nil
		},
	}
)

func init() {
	backupCmd.Flags().BoolVarP(&listBackups, "list", "l", false, "List available backups")
}
