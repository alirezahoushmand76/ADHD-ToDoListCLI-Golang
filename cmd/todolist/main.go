package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/user/todolist/cmd/todolist/cmd"
	"github.com/user/todolist/internal/client"
	"github.com/user/todolist/internal/ui"
)

var (
	serverAddr string
)

func init() {
	flag.StringVar(&serverAddr, "server", "localhost:8080", "Address of the TodoList server")
}

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

	// Parse flags before cobra gets to them
	flag.Parse()

	// Remove the parsed flags from os.Args
	os.Args = append(os.Args[:1], flag.Args()...)

	// Connect to the server if specified
	if serverAddr != "" {
		todoClient, err := client.NewClient(serverAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to local mode\n")
		} else {
			defer todoClient.Close()
			cmd.InitializeWithClient(todoClient)
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
