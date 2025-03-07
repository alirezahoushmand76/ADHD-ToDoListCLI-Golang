# TodoList - ADHD-Friendly CLI Task Manager

TodoList is a command-line task management application designed specifically for users with ADHD. It provides a distraction-free interface with features that help maintain focus and manage tasks effectively.

## Features

- **Task Management**: Add, list, complete, and delete tasks with ease
- **Priority Levels**: Assign low, medium, or high priority to tasks
- **Categories**: Organize tasks into different categories (work, personal, urgent, etc.)
- **Due Dates & Reminders**: Set due dates and reminders for tasks
- **Brain Dump Mode**: Quickly add multiple tasks without interruption
- **Focus Mode**: Get suggestions for the next task to work on based on priority and urgency
- **Pomodoro Timer**: Built-in timer for focused work sessions
- **Data Persistence**: Tasks are stored locally in JSON format
- **Backup & Restore**: Create backups of your tasks and restore them when needed
- **Colorful Output**: Visual cues make tasks more readable and engaging

## Installation

### Prerequisites

- Go 1.22 or later

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/user/todolist.git
   cd todolist
   ```

2. Build the application:
   ```
   make build
   ```

3. (Optional) Install the binary to your system:
   ```
   make install
   ```

### Using Docker

#### Prerequisites

- Docker
- Docker Compose

#### Running with Docker Compose

1. Build and run the application:
   ```
   docker-compose build
   docker-compose run todolist
   ```

2. Or use the provided shell script:
   ```
   ./todolist.sh
   ```

3. Run commands:
   ```
   ./todolist.sh add "Complete project report" --priority high
   ./todolist.sh list
   ```

## Usage

### Basic Commands

- **Add a task**:
  ```
  todolist add "Complete project report" --priority high --category work --due "2023-12-31"
  ```

- **List tasks**:
  ```
  todolist list
  todolist list --category work
  todolist list --priority high
  todolist list --all
  ```

- **Complete a task**:
  ```
  todolist complete [task_id]
  ```

- **Delete a task**:
  ```
  todolist delete [task_id]
  ```

### ADHD-Specific Features

- **Brain Dump Mode**:
  ```
  todolist dump
  ```

- **Focus Mode**:
  ```
  todolist focus
  ```

- **Pomodoro Timer**:
  ```
  todolist pomodoro [task_id]
  todolist pomodoro [task_id] --duration 30
  ```

### Data Management

- **Backup tasks**:
  ```
  todolist backup
  todolist backup --list
  ```

- **Restore tasks**:
  ```
  todolist restore [backup_file_or_index]
  ```

## Makefile Commands

The project includes a Makefile for common operations:

- `make build`: Build the application
- `make run`: Build and run the application
- `make clean`: Clean build artifacts
- `make test`: Run tests
- `make deps`: Update dependencies
- `make install`: Install the binary to your system

## Docker Deployment

### Local Deployment

1. Build the Docker image:
   ```
   docker-compose build
   ```

2. Run the application:
   ```
   docker-compose run todolist [command]
   ```

3. Use the shell script wrapper for convenience:
   ```
   ./todolist.sh [command]
   ```

### Remote Deployment

1. Build and tag the Docker image:
   ```
   docker build -t yourusername/todolist:latest .
   ```

2. Push the image to a Docker registry:
   ```
   docker push yourusername/todolist:latest
   ```

3. On the remote server, pull and run the image:
   ```
   docker pull yourusername/todolist:latest
   docker run -v todolist_data:/home/appuser/.todolist yourusername/todolist:latest [command]
   ```

## CI/CD Pipeline

This project includes a GitHub Actions CI/CD pipeline that automates testing, building, and deployment:

### Pipeline Stages

1. **Test**: Runs unit tests to ensure code quality
2. **Build**: Builds and pushes the Docker image to Docker Hub
3. **Deploy**: Deploys the application to a server (on main/master branch)
4. **Release**: Creates GitHub releases with binaries for multiple platforms (on tags)

### Setting Up GitHub Secrets

To use the CI/CD pipeline, you need to set up the following GitHub secrets:

- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token
- `SSH_PRIVATE_KEY`: SSH private key for server access
- `SSH_USER`: SSH username for server access
- `SERVER_IP`: IP address of your deployment server

### Manual Deployment

You can also deploy manually using the provided deployment script:

1. Update the Docker image name in `deploy.sh`:
   ```bash
   DOCKER_IMAGE="yourusername/todolist:latest"
   ```

2. Run the deployment script on your server:
   ```bash
   ./deploy.sh
   ```

## Project Structure

- `/cmd/todolist`: CLI application entry point and commands
- `/internal/app`: Core application logic
- `/internal/models`: Data models
- `/internal/storage`: Data persistence
- `/internal/ui`: User interface utilities
- `/internal/utils`: Utility functions

## Future Enhancements

- AI-assisted task prioritization
- Natural language processing for task entry
- Sync with web and mobile versions
- Integration with calendar applications

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 