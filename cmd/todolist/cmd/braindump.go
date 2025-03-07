package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/user/todolist/internal/ui"
)

var brainDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Enter brain dump mode",
	Long:  `Enter brain dump mode to quickly add multiple tasks without interruption.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set up signal handling for graceful exit with Ctrl+C
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Println("\nBrain dump mode interrupted. Tasks entered so far have been saved.")
			os.Exit(0)
		}()

		ui.PrintInfo("🧠 BRAIN DUMP MODE")
		ui.PrintInfo("Quickly add tasks without interruption. Enter one task per line.")
		ui.PrintInfo("Leave a line empty when you're done.")
		ui.PrintInfo("Press Ctrl+C at any time to exit and save tasks entered so far.")
		fmt.Println()

		err := todoApp.BrainDump()
		if err != nil {
			return fmt.Errorf("brain dump mode error: %w", err)
		}

		return nil
	},
	Example: `  todolist dump`,
}
