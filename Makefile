
ROOT_DIR := D:/GIN/real-time-bus-tracking-service
BUILD_DIR := $(ROOT_DIR)/build/bin
TMP_DIR := $(ROOT_DIR)/build/tmp
SERVICES := api-gateway bus-service

$(shell mkdir -p $(BUILD_DIR) $(TMP_DIR)/logs)

.PHONY: all build dev dev-% prod prod-% docker docker-down clean deps test lint

all: build

build:
	@for service in $(SERVICES); do \
		go build -o $(BUILD_DIR)/$$service.exe $(ROOT_DIR)/cmd/$$service/main.go; \
	done

dev:
	@echo "Starting Development environment"
	@for service in $(SERVICES); do \
		echo "Starting $$service in development mode..."; \
		(cd $(ROOT_DIR)/cmd/$$service && air &); \
	done

dev-%:
	@echo "Starting $* in development mode..."
	@cd $(ROOT_DIR)/cmd/$* && air

prod:
	@for service in $(SERVICES); do \
		echo "Building and starting $$service in production mode..."; \
		go build -o $(BUILD_DIR)/$$service.exe $(ROOT_DIR)/cmd/$$service/main.go; \
		$(BUILD_DIR)/$$service.exe & \
	done

prod-%:
	@echo "Building and starting $* in production mode..."
	@go build -o $(BUILD_DIR)/$*.exe $(ROOT_DIR)/cmd/$*/main.go
	@$(BUILD_DIR)/$*.exe &

docker:
	@echo "Starting Docker Compose..."
	@cd $(ROOT_DIR)/config && docker-compose up --build -d

docker-down:
	@echo "Stopping Docker Compose..."
	@cd $(ROOT_DIR)/config && docker-compose down

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/*.exe $(TMP_DIR)/logs/*.log

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/air-verse/air@latest

test:
	@echo "Running tests..."
	@go test -v ./tests/...

lint:
	@echo "Running linter..."
	@go fmt ./...
	@go vet ./...
