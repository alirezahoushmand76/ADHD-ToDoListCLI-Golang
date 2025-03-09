package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/user/todolist/internal/app"
	"github.com/user/todolist/internal/models"
	"github.com/user/todolist/internal/protocol"
	"github.com/user/todolist/internal/storage"
)

var (
	port    = flag.String("port", "8080", "Port to listen on")
	dataDir = flag.String("data-dir", "", "Data directory (defaults to ~/.todolist)")
)

func main() {
	flag.Parse()

	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetPrefix("[TodoList Server] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize app configuration
	config := app.DefaultConfig()
	if *dataDir != "" {
		config.DataDir = *dataDir
		config.StorageFile = filepath.Join(*dataDir, "tasks.json")
		config.BackupDir = filepath.Join(*dataDir, "backups")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	if err := os.MkdirAll(config.BackupDir, 0755); err != nil {
		log.Fatalf("Failed to create backup directory: %v", err)
	}

	// Initialize storage
	store := storage.NewJSONStorage(config.StorageFile)
	if err := store.Initialize(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize app
	todoApp, err := app.NewApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start TCP server
	addr := fmt.Sprintf(":%s", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Check if it's a port conflict
		if strings.Contains(err.Error(), "address already in use") {
			log.Fatalf("Port %s is already in use. Please try a different port with --port flag or kill the process using this port.", *port)
		}
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("TodoList server started on %s", addr)

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Shutting down server...")
		listener.Close()
		os.Exit(0)
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn, todoApp)
	}
}

func handleConnection(conn net.Conn, todoApp *app.App) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", clientAddr)

	reader := bufio.NewReader(conn)

	for {
		// Read request
		requestData, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client %s disconnected", clientAddr)
			} else {
				log.Printf("Error reading from client %s: %v", clientAddr, err)
			}
			return
		}

		// Parse request
		var request protocol.Request
		if err := json.Unmarshal(requestData, &request); err != nil {
			sendErrorResponse(conn, fmt.Sprintf("Invalid request format: %v", err))
			continue
		}

		// Process request
		response := processRequest(todoApp, request)

		// Send response
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}
		responseData = append(responseData, '\n')

		if _, err := conn.Write(responseData); err != nil {
			log.Printf("Error sending response to client %s: %v", clientAddr, err)
			return
		}
	}
}

func processRequest(todoApp *app.App, request protocol.Request) protocol.Response {
	var response protocol.Response

	switch request.Operation {
	case protocol.OpAddTask:
		var addReq protocol.AddTaskRequest
		if err := json.Unmarshal(request.Payload, &addReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid add task request: %v", err))
		}

		task := models.NewTask(
			addReq.Title,
			addReq.Description,
			addReq.Priority,
			addReq.Category,
			addReq.DueDate,
			addReq.ReminderAt,
		)

		if err := todoApp.AddTask(
			addReq.Title,
			addReq.Description,
			addReq.Priority,
			addReq.Category,
			addReq.DueDate,
			addReq.ReminderAt,
		); err != nil {
			return errorResponse(fmt.Sprintf("Failed to add task: %v", err))
		}

		taskResp := protocol.TaskResponse{Task: task}
		payload, _ := json.Marshal(taskResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpGetTask:
		var idReq protocol.IDRequest
		if err := json.Unmarshal(request.Payload, &idReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid get task request: %v", err))
		}

		task, err := todoApp.GetTask(idReq.ID)
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to get task: %v", err))
		}

		taskResp := protocol.TaskResponse{Task: task}
		payload, _ := json.Marshal(taskResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpGetAllTasks:
		tasks, err := todoApp.GetAllTasks()
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to get tasks: %v", err))
		}

		tasksResp := protocol.TasksResponse{Tasks: tasks}
		payload, _ := json.Marshal(tasksResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpGetTasksByCategory:
		var catReq protocol.CategoryRequest
		if err := json.Unmarshal(request.Payload, &catReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid category request: %v", err))
		}

		tasks, err := todoApp.GetTasksByCategory(catReq.Category)
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to get tasks by category: %v", err))
		}

		tasksResp := protocol.TasksResponse{Tasks: tasks}
		payload, _ := json.Marshal(tasksResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpGetTasksByPriority:
		var prioReq protocol.PriorityRequest
		if err := json.Unmarshal(request.Payload, &prioReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid priority request: %v", err))
		}

		tasks, err := todoApp.GetTasksByPriority(prioReq.Priority)
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to get tasks by priority: %v", err))
		}

		tasksResp := protocol.TasksResponse{Tasks: tasks}
		payload, _ := json.Marshal(tasksResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpUpdateTask:
		var taskReq protocol.TaskResponse
		if err := json.Unmarshal(request.Payload, &taskReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid update task request: %v", err))
		}

		if err := todoApp.UpdateTask(taskReq.Task); err != nil {
			return errorResponse(fmt.Sprintf("Failed to update task: %v", err))
		}

		response.Success = true

	case protocol.OpDeleteTask:
		var idReq protocol.IDRequest
		if err := json.Unmarshal(request.Payload, &idReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid delete task request: %v", err))
		}

		if err := todoApp.DeleteTask(idReq.ID); err != nil {
			return errorResponse(fmt.Sprintf("Failed to delete task: %v", err))
		}

		response.Success = true

	case protocol.OpCompleteTask:
		var idReq protocol.IDRequest
		if err := json.Unmarshal(request.Payload, &idReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid complete task request: %v", err))
		}

		if err := todoApp.CompleteTask(idReq.ID); err != nil {
			return errorResponse(fmt.Sprintf("Failed to complete task: %v", err))
		}

		response.Success = true

	case protocol.OpBackup:
		var backupReq protocol.BackupRequest
		if err := json.Unmarshal(request.Payload, &backupReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid backup request: %v", err))
		}

		filename, err := todoApp.BackupTasks()
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to backup tasks: %v", err))
		}

		backupResp := protocol.BackupResponse{Filename: filename}
		payload, _ := json.Marshal(backupResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpRestore:
		var restoreReq protocol.RestoreRequest
		if err := json.Unmarshal(request.Payload, &restoreReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid restore request: %v", err))
		}

		if err := todoApp.RestoreTasks(restoreReq.Filename); err != nil {
			return errorResponse(fmt.Sprintf("Failed to restore tasks: %v", err))
		}

		response.Success = true

	case protocol.OpListBackups:
		backups, err := todoApp.ListBackups()
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to list backups: %v", err))
		}

		backupsResp := protocol.ListBackupsResponse{Backups: backups}
		payload, _ := json.Marshal(backupsResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpBrainDump:
		if err := todoApp.BrainDump(); err != nil {
			return errorResponse(fmt.Sprintf("Failed to perform brain dump: %v", err))
		}

		response.Success = true

	case protocol.OpFocusMode:
		task, err := todoApp.FocusMode()
		if err != nil {
			return errorResponse(fmt.Sprintf("Failed to enter focus mode: %v", err))
		}

		taskResp := protocol.TaskResponse{Task: task}
		payload, _ := json.Marshal(taskResp)
		response.Success = true
		response.Payload = payload

	case protocol.OpStartPomodoro:
		var pomReq protocol.PomodoroRequest
		if err := json.Unmarshal(request.Payload, &pomReq); err != nil {
			return errorResponse(fmt.Sprintf("Invalid pomodoro request: %v", err))
		}

		if err := todoApp.StartPomodoro(pomReq.TaskID, pomReq.CustomDuration); err != nil {
			return errorResponse(fmt.Sprintf("Failed to start pomodoro: %v", err))
		}

		response.Success = true

	default:
		return errorResponse(fmt.Sprintf("Unknown operation: %s", request.Operation))
	}

	return response
}

func errorResponse(message string) protocol.Response {
	return protocol.Response{
		Success: false,
		Error:   message,
	}
}

func sendErrorResponse(conn net.Conn, message string) {
	response := errorResponse(message)
	responseData, _ := json.Marshal(response)
	responseData = append(responseData, '\n')
	conn.Write(responseData)
}
