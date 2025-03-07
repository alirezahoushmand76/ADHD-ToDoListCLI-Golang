# ADHD-Friendly To-Do List CLI Application Documentation

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage Guide](#usage-guide)
5. [Architecture](#architecture)
6. [Development Guide](#development-guide)
7. [Deployment](#deployment)
8. [Troubleshooting](#troubleshooting)
9. [Contributing](#contributing)
10. [License](#license)

## Introduction

The ADHD-Friendly To-Do List CLI is a command-line task management application designed specifically for users with ADHD. It provides a distraction-free interface with features that help maintain focus and manage tasks effectively.

This application is built with Go and follows best practices for CLI application development. It includes features specifically designed to help users with ADHD, such as Brain Dump mode, Focus Mode, and Pomodoro Timer integration.

### Key Benefits

- **Distraction-free interface**: Simple CLI interface without visual clutter
- **ADHD-specific features**: Tools designed to help with focus and task management
- **Efficient workflow**: Quick commands for common operations
- **Data persistence**: Tasks are stored locally and can be backed up
- **Visual cues**: Colorful output for better task differentiation

## Features

### Task Management

- **Add tasks**: Create new tasks with title, description, priority, category, due date, and reminder
- **List tasks**: View tasks with filtering options (by category, priority, completion status)
- **Complete tasks**: Mark tasks as completed
- **Delete tasks**: Remove tasks from the list
- **Categories**: Organize tasks into groups (work, personal, urgent, etc.)
- **Priorities**: Assign low, medium, or high priority to tasks
- **Due dates**: Set deadlines for tasks
- **Reminders**: Set reminder times for tasks

### ADHD-Specific Features

- **Brain Dump Mode**: Quickly add multiple tasks without interruption
- **Focus Mode**: Get suggestions for the next task based on priority, urgency, and other factors
- **Pomodoro Timer**: Built-in timer for focused work sessions

### Data Management

- **JSON Storage**: Tasks are stored locally in JSON format
- **Backup & Restore**: Create backups of your tasks and restore them when needed

## Installation

### Prerequisites

- Go 1.22 or later
- Docker (optional, for containerized usage)

### Option 1: Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/alirezahoushmand76/ADHD-ToDoListCLI-Golang.git
   cd ADHD-ToDoListCLI-Golang
   ```

2. Build the application:
   ```bash
   make build
   ```

3. (Optional) Install the binary to your system:
   ```bash
   make install
   ```

### Option 2: Using Docker

1. Pull the Docker image:
   ```bash
   docker pull alirezahoushmand76/todolist:latest
   ```

2. Create a volume for data persistence:
   ```bash
   docker volume create todolist_data
   ```

3. Create a wrapper script:
   ```bash
   echo '#!/bin/bash
   docker run --rm -v todolist_data:/home/appuser/.todolist alirezahoushmand76/todolist:latest "$@"' > todolist
   chmod +x todolist
   sudo mv todolist /usr/local/bin/
   ```

### Option 3: Download Binary Releases

1. Go to the [Releases page](https://github.com/alirezahoushmand76/ADHD-ToDoListCLI-Golang/releases)
2. Download the appropriate binary for your platform (Linux, macOS, or Windows)
3. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x todolist-linux-amd64  # or todolist-darwin-amd64
   ```
4. Move the binary to a location in your PATH:
   ```bash
   sudo mv todolist-linux-amd64 /usr/local/bin/todolist
   ```

## Usage Guide

### Basic Commands

#### Adding Tasks

```bash
# Basic task
todolist add "Complete project report"

# Task with priority and category
todolist add "Complete project report" --priority high --category work

# Task with due date and reminder
todolist add "Call doctor" --due tomorrow --reminder "tomorrow 9:00"

# Task with description
todolist add "Research topic" --description "Find sources for the research paper" --category work
```

#### Listing Tasks

```bash
# List all incomplete tasks
todolist list

# List all tasks (including completed)
todolist list --all

# List tasks by category
todolist list --category work

# List tasks by priority
todolist list --priority high
```

#### Completing Tasks

```bash
# Mark a task as completed
todolist complete 1741359296120413000  # Replace with actual task ID
```

#### Deleting Tasks

```bash
# Delete a task
todolist delete 1741359296120413000  # Replace with actual task ID

# Delete a task without confirmation
todolist delete 1741359296120413000 --force
```

### ADHD-Specific Features

#### Brain Dump Mode

Brain Dump mode allows you to quickly add multiple tasks without interruption:

```bash
todolist dump
```

In Brain Dump mode:
- Enter task titles one per line
- Leave a line empty to finish
- Type 'q', 'quit', or 'exit' to exit
- Type 'help' for help
- Press Ctrl+C to exit at any time

#### Focus Mode

Focus mode suggests the next task to work on based on priority, urgency, and other factors:

```bash
todolist focus
```

#### Pomodoro Timer

The Pomodoro Timer helps you focus on a task for a set period of time:

```bash
# Start a Pomodoro timer for a task with default duration (25 minutes)
todolist pomodoro 1741359296120413000  # Replace with actual task ID

# Start a Pomodoro timer with custom duration
todolist pomodoro 1741359296120413000 --duration 30  # 30-minute work session
```

### Data Management

#### Backup

```bash
# Create a backup
todolist backup

# List available backups
todolist backup --list
```

#### Restore

```bash
# Restore from a backup by index
todolist restore 1

# Restore from a backup file
todolist restore /path/to/backup/file.json

# Restore without confirmation
todolist restore 1 --force
```

## Architecture

The application follows a modular architecture with clear separation of concerns:

### Project Structure

```
todolist-in-golang/
├── cmd/
│   └── todolist/
│       ├── cmd/
│       │   ├── add.go
│       │   ├── backup.go
│       │   ├── braindump.go
│       │   ├── complete.go
│       │   ├── delete.go
│       │   ├── focus.go
│       │   ├── list.go
│       │   ├── pomodoro.go
│       │   ├── restore.go
│       │   └── root.go
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── models/
│   │   └── task.go
│   ├── storage/
│   │   ├── json_storage.go
│   │   └── storage.go
│   ├── ui/
│   │   └── ui.go
│   └── utils/
│       ├── focus.go
│       └── pomodoro.go
├── .github/
│   └── workflows/
│       └── ci-cd.yml
├── Dockerfile
├── docker-compose.yml
├── deploy.sh
├── todolist.sh
├── Makefile
└── README.md
```

### Components

1. **cmd/todolist**: CLI application entry point and commands
   - **main.go**: Application entry point with global error handling
   - **cmd/**: Command implementations using Cobra

2. **internal/app**: Core application logic
   - **app.go**: Application core that ties together storage and business logic

3. **internal/models**: Data models
   - **task.go**: Task model with fields and methods

4. **internal/storage**: Data persistence
   - **storage.go**: Storage interface
   - **json_storage.go**: JSON file-based storage implementation

5. **internal/ui**: User interface utilities
   - **ui.go**: UI utilities for colorful output

6. **internal/utils**: Utility functions
   - **focus.go**: Focus mode implementation
   - **pomodoro.go**: Pomodoro timer implementation

## Development Guide

### Setting Up Development Environment

1. Install Go 1.22 or later
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```

### Building the Application

```bash
make build
```

### Running Tests

```bash
make test
```

### Code Style and Conventions

- Follow Go's official [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format code
- Document exported functions and types
- Handle errors explicitly
- Use meaningful variable and function names

### Adding New Features

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Implement your feature

3. Add tests

4. Submit a pull request

## Deployment

### Docker Deployment

#### Local Deployment

1. Build the Docker image:
   ```bash
   docker-compose build
   ```

2. Run the application:
   ```bash
   docker-compose run todolist [command]
   ```

3. Use the shell script wrapper for convenience:
   ```bash
   ./todolist.sh [command]
   ```

#### Remote Deployment

1. Build and tag the Docker image:
   ```bash
   docker build -t yourusername/todolist:latest .
   ```

2. Push the image to a Docker registry:
   ```bash
   docker push yourusername/todolist:latest
   ```

3. On the remote server, pull and run the image:
   ```bash
   docker pull yourusername/todolist:latest
   docker run -v todolist_data:/home/appuser/.todolist yourusername/todolist:latest [command]
   ```

### CI/CD Pipeline

This project includes a GitHub Actions CI/CD pipeline that automates testing, building, and deployment:

#### Pipeline Stages

1. **Test**: Runs unit tests to ensure code quality
2. **Build**: Builds and pushes the Docker image to Docker Hub
3. **Deploy**: Deploys the application to a server (on main/master branch)
4. **Release**: Creates GitHub releases with binaries for multiple platforms (on tags)

#### Setting Up GitHub Secrets

To use the CI/CD pipeline, you need to set up the following GitHub secrets:

- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token
- `SSH_PRIVATE_KEY`: SSH private key for server access
- `SSH_USER`: SSH username for server access
- `SERVER_IP`: IP address of your deployment server

#### Manual Deployment

You can also deploy manually using the provided deployment script:

1. Update the Docker image name in `deploy.sh`:
   ```bash
   DOCKER_IMAGE="yourusername/todolist:latest"
   ```

2. Run the deployment script on your server:
   ```bash
   ./deploy.sh
   ```

## Troubleshooting

### Common Issues and Solutions

#### Application Crashes

If the application crashes, check the error message. The application includes global error handling to provide helpful error messages.

#### Task ID Format

When using task IDs, make sure to use the exact ID without square brackets. For example:
```bash
# Correct
todolist complete 1741359296120413000

# Incorrect
todolist complete [1741359296120413000]
```

#### Docker Issues

If you encounter issues with Docker:

1. Make sure Docker is running:
   ```bash
   docker info
   ```

2. Check if the image exists:
   ```bash
   docker images | grep todolist
   ```

3. Check if the volume exists:
   ```bash
   docker volume ls | grep todolist_data
   ```

#### CI/CD Pipeline Issues

If the CI/CD pipeline fails:

1. Check the GitHub Actions logs for errors
2. Verify that all required secrets are set up correctly
3. Make sure the Docker Hub credentials are valid

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes
4. Add tests for your changes
5. Run the tests:
   ```bash
   go test -v ./...
   ```
6. Commit your changes:
   ```bash
   git commit -m "Add your feature"
   ```
7. Push to the branch:
   ```bash
   git push origin feature/your-feature-name
   ```
8. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Appendix: Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `add` | Add a new task | `todolist add "Complete project report" --priority high` |
| `list` | List tasks | `todolist list --category work` |
| `complete` | Mark a task as completed | `todolist complete 1741359296120413000` |
| `delete` | Delete a task | `todolist delete 1741359296120413000` |
| `dump` | Enter brain dump mode | `todolist dump` |
| `focus` | Enter focus mode | `todolist focus` |
| `pomodoro` | Start a Pomodoro timer | `todolist pomodoro 1741359296120413000 --duration 30` |
| `backup` | Create or list backups | `todolist backup --list` |
| `restore` | Restore from a backup | `todolist restore 1` |
| `version` | Show version information | `todolist version` | 