.PHONY: build run clean test

# Build variables
BINARY_NAME=todolist
BUILD_DIR=build

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

all: clean build

# Build the client and server binaries
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-client ./cmd/todolist
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-server ./cmd/server

# Build the client binary
build-client:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-client ./cmd/todolist

# Build the server binary
build-server:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-server ./cmd/server

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) tidy

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Specific commands
add:
	./$(BUILD_DIR)/$(BINARY_NAME) add $(filter-out $@,$(MAKECMDGOALS))

list:
	./$(BUILD_DIR)/$(BINARY_NAME) list $(filter-out $@,$(MAKECMDGOALS))

complete:
	./$(BUILD_DIR)/$(BINARY_NAME) complete $(filter-out $@,$(MAKECMDGOALS))

delete:
	./$(BUILD_DIR)/$(BINARY_NAME) delete $(filter-out $@,$(MAKECMDGOALS))

dump:
	./$(BUILD_DIR)/$(BINARY_NAME) dump

focus:
	./$(BUILD_DIR)/$(BINARY_NAME) focus

pomodoro:
	./$(BUILD_DIR)/$(BINARY_NAME) pomodoro $(filter-out $@,$(MAKECMDGOALS))

backup:
	./$(BUILD_DIR)/$(BINARY_NAME) backup $(filter-out $@,$(MAKECMDGOALS))

restore:
	./$(BUILD_DIR)/$(BINARY_NAME) restore $(filter-out $@,$(MAKECMDGOALS))

# Run the server
run-server: build-server
	./$(BUILD_DIR)/$(BINARY_NAME)-server

# Run the client
run-client: build-client
	./$(BUILD_DIR)/$(BINARY_NAME)-client

# Allow passing arguments to commands
%:
	@: 