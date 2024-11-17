GO=go
GOFMT=gofmt
GOFLAGS=-mod=readonly
PACKAGE=github.com/JBSWE/load-balancer

SRCDIR=cmd/server
INTERNALDIR=internal

BINARY_NAME=load-balancer

BUILD_FLAGS=-o $(BINARY_NAME)

GO_TEST = $(GO) test
INTEGRATION_DIR = ./integration
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_FILE = docker-compose.yml
DOCKER_NETWORK = load-balancer-network

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building the Go project..."
	$(GO) build -o $(BINARY_NAME) ./cmd/server

.PHONY: run
run: build
	@echo "Running the application..."
	./$(BINARY_NAME) ./cmd/server

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v $(shell go list ./... | grep -v '/integration')

.PHONY: fmt
fmt:
	@echo "Formatting Go files..."
	$(GOFMT) -s -w .

.PHONY: vet
vet:
	@echo "Running go vet on the Go files..."
	$(GO) vet ./...

.PHONY: clean
clean:
	@echo "Cleaning the project..."
	rm -f $(BINARY_NAME)

.PHONY: tidy
tidy:
	@echo "Tidying up Go modules..."
	$(GO) mod tidy

run-verbose: build
	@echo "Running the server with verbose output..."
	./$(BINARY_NAME) -v

all: run-integration-tests

build:
	$(GO) build -o main ./cmd/server

up:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

down:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

run-integration-tests: up
	@echo "Running integration tests..."
	$(GO_TEST) $(INTEGRATION_DIR) -v

clean:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v
	rm -rf bin

build-and-test: build up run-integration-tests down

rebuild: down build up run-integration-tests down

.PHONY: build up down run-integration-tests clean rebuild build-and-test
