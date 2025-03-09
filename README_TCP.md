# TodoList TCP Client-Server Architecture

This version of the TodoList application has been refactored to use a client-server architecture with TCP networking. This allows multiple clients to connect to a central server that manages the task data.

## Architecture

The application is split into two main components:

1. **Server**: Manages the storage of tasks and handles client requests
2. **Client**: Connects to the server and provides a command-line interface for users

## Running the Server

To start the server, run:

```bash
go run cmd/server/main.go
```

By default, the server listens on port 8080. You can change this with the `--port` flag:

```bash
go run cmd/server/main.go --port 9090
```

The server stores data in the `~/.todolist` directory by default. You can change this with the `--data-dir` flag:

```bash
go run cmd/server/main.go --data-dir /path/to/data
```

## Running the Client

To start the client, run:

```bash
go run cmd/todolist/main.go
```

By default, the client will run in standalone mode. To connect to a server, use the `-server` flag:

```bash
go run cmd/todolist/main.go -server localhost:8080
```

You can specify a different server address:

```bash
go run cmd/todolist/main.go -server 192.168.1.100:9090
```

**Note**: The `-server` flag must be specified before any subcommands:

```bash
# Correct
go run cmd/todolist/main.go -server localhost:8080 list

# Incorrect
go run cmd/todolist/main.go list -server localhost:8080
```

## Building the Application

To build both the client and server, run:

```bash
go build -o todolist-server cmd/server/main.go
go build -o todolist-client cmd/todolist/main.go
```

## Protocol

The client and server communicate using a JSON-based protocol. Each message consists of:

1. An operation type (e.g., `ADD_TASK`, `GET_TASK`, etc.)
2. A payload specific to the operation
3. A newline character (`\n`) to delimit messages

## Features

The client-server architecture supports all the features of the original TodoList application:

- Adding, listing, completing, and deleting tasks
- Filtering tasks by category or priority
- Setting due dates and reminders
- Brain dump mode
- Focus mode
- Pomodoro timer

## Benefits of the Client-Server Architecture

- **Centralized Data Storage**: All tasks are stored on the server, allowing multiple clients to access the same data.
- **Remote Access**: Clients can connect to the server from different machines.
- **Scalability**: The server can handle multiple client connections simultaneously.
- **Separation of Concerns**: The client handles user interaction, while the server manages data storage and business logic.

## Security Considerations

This implementation uses plain TCP without encryption. For production use, consider adding TLS encryption and authentication mechanisms.

## Troubleshooting

### "Address already in use" error

If you see an error like this when starting the server:

```
Failed to start server: listen tcp :8080: bind: address already in use
```

It means that port 8080 is already being used by another process. You have two options:

1. **Kill the existing process**:
   ```bash
   # Find the process using port 8080
   lsof -i :8080
   
   # Kill the process (replace PID with the actual process ID)
   kill PID
   ```

2. **Use a different port**:
   ```bash
   # Start the server on a different port
   go run cmd/server/main.go --port 9090
   
   # Connect the client to the server on the new port
   go run cmd/todolist/main.go -server localhost:9090 list
   ```

### Client can't connect to server

If the client can't connect to the server, make sure:

1. The server is running
2. You're using the correct server address and port
3. There are no firewall rules blocking the connection
4. The `-server` flag is specified before any subcommands 