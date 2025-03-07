package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/app"
	"github.com/user/todolist/internal/ui"
)

var (
	todoApp *app.App
	rootCmd = &cobra.Command{
		Use:   "todolist",
		Short: "An ADHD-friendly To-Do List CLI application",
		Long: `
TodoList is an ADHD-friendly To-Do List CLI application designed to help users 
with ADHD manage their tasks effectively. It includes features like task management, 
brain dump mode, focus mode, and Pomodoro timer integration.
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip initialization for completion and help commands
			if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "version" {
				return nil
			}

			var err error
			todoApp, err = app.NewApp(nil)
			if err != nil {
				return fmt.Errorf("failed to initialize application: %w", err)
			}
			return nil
		},
		// Add a global error handler for all commands
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(brainDumpCmd)
	rootCmd.AddCommand(focusCmd)
	rootCmd.AddCommand(pomodoroCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TodoList v1.0.0")
		},
	})

	// Set up custom error handling for all commands
	cobra.OnInitialize(func() {
		for _, cmd := range rootCmd.Commands() {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
		}
	})
}

// handleError handles errors in a user-friendly way
func handleError(err error) {
	if err != nil {
		// Extract the error message without the usage information
		errMsg := err.Error()
		if idx := strings.Index(errMsg, "\nUsage:"); idx != -1 {
			errMsg = errMsg[:idx]
		}

		ui.PrintError("Error: %v", errMsg)
		os.Exit(1)
	}
}
