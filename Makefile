GO=go
GOFMT=gofmt
GOFLAGS=-mod=readonly
PACKAGE=github.com/JBSWE/load-balancer

SRCDIR=cmd/server
INTERNALDIR=internal

BINARY_NAME=load-balancer

BUILD_FLAGS=-o $(BINARY_NAME)

.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the Go project..."
	$(GO) build -o $(BINARY_NAME) ./cmd/server

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	./$(BINARY_NAME) ./cmd/server

# Test the application unit tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Format Go files (using gofmt)
.PHONY: fmt
fmt:
	@echo "Formatting Go files..."
	$(GOFMT) -s -w .

# Run vet for basic linting
.PHONY: vet
vet:
	@echo "Running go vet on the Go files..."
	$(GO) vet ./...

# Clean the build artifacts
.PHONY: clean
clean:
	@echo "Cleaning the project..."
	rm -f $(BINARY_NAME)

# Install dependencies
.PHONY: tidy
tidy:
	@echo "Tidying up Go modules..."
	$(GO) mod tidy

# Run the application with `make run` (combined build and run)
run-verbose: build
	@echo "Running the server with verbose output..."
	./$(BINARY_NAME) -v
