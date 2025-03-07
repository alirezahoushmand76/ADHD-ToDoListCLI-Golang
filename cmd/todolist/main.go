package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/user/todolist/cmd/todolist/cmd"
	"github.com/user/todolist/internal/ui"
)

func main() {
	// Set up panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "")
			ui.PrintError("The application encountered an unexpected error")
			ui.PrintError("Error details: %v", r)
			ui.PrintError("This is likely a bug. Please report it with the following stack trace:")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, string(debug.Stack()))
			os.Exit(1)
		}
	}()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
