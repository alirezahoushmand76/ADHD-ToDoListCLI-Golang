package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/app"
	"github.com/user/todolist/internal/client"
	"github.com/user/todolist/internal/storage"
	"github.com/user/todolist/internal/ui"
)

var (
	todoClient *client.Client
	dataDir    string
	verbose    bool
	todoApp    *app.App
)

var rootCmd = &cobra.Command{
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

		// If we're using the client, we don't need to initialize the app
		if todoClient != nil {
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

// Execute executes the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return handleError(err)
	}
	return nil
}

// InitializeWithClient initializes the command with a client
func InitializeWithClient(client *client.Client) {
	todoClient = client
}

func init() {
	// Define persistent flags for the root command
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Data directory (defaults to ~/.todolist)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// If no client is provided, we're running in standalone mode (for backward compatibility)
	cobra.OnInitialize(initConfig)

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

// initConfig initializes the application configuration
func initConfig() {
	// Skip initialization if client is already set
	if todoClient != nil {
		return
	}

	// Legacy standalone mode - initialize app directly
	config := app.DefaultConfig()
	if dataDir != "" {
		config.DataDir = dataDir
		config.StorageFile = filepath.Join(dataDir, "tasks.json")
		config.BackupDir = filepath.Join(dataDir, "backups")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating backup directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	store := storage.NewJSONStorage(config.StorageFile)
	if err := store.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize app
	var err error
	todoApp, err = app.NewApp(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing app: %v\n", err)
		os.Exit(1)
	}
}

// handleError handles errors in a user-friendly way
func handleError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's a known error type
	switch {
	case strings.Contains(err.Error(), "not found"):
		ui.PrintError("Task not found. Please check the ID and try again.")
	case strings.Contains(err.Error(), "invalid input"):
		ui.PrintError("Invalid input: %v", err)
	default:
		ui.PrintError("Error: %v", err)
	}

	return err
}
