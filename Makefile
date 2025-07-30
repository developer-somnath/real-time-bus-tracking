
# Define variables
ROOT_DIR := D:/GIN/real-time-bus-tracking-service
BUILD_DIR := $(ROOT_DIR)/build/bin
TMP_DIR := $(ROOT_DIR)/build/tmp

# Ensure directories exist
$(shell mkdir -p $(BUILD_DIR) $(TMP_DIR)/logs)

# Services
SERVICES := api-gateway bus-service route-service trip-service user-service notification-service analytics-service payment-service geolocation-service driver-service config-service event-service

# Default target
.PHONY: all
all: build

# Build all services
.PHONY: build
build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		go build -o $(BUILD_DIR)/$$service.exe $(ROOT_DIR)/cmd/$$service/main.go; \
	done

# Run Air for a specific service
.PHONY: air-%
air-%:
	@echo "Starting Air for $*..."
	@cd $(ROOT_DIR)/cmd/$* && air

# Run all services with Air
.PHONY: air-all
air-all:
	@echo "Starting all services with Air..."
	@for service in $(SERVICES); do \
		(cd $(ROOT_DIR)/cmd/$$service && air &); \
	done

# Run Docker Compose
.PHONY: docker
docker:
	@echo "Starting Docker Compose..."
	@cd $(ROOT_DIR)/config && docker-compose up --build

# Stop Docker Compose
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose..."
	@cd $(ROOT_DIR)/config && docker-compose down

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./tests/...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/*.exe
	@rm -rf $(TMP_DIR)/logs/*.log

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/air-verse/air@latest
```

### Explanation of Targets
- **all**: Default target, runs `build`.
- **build**: Compiles all services (e.g., `api-gateway`, `bus-service`) into `build/bin/<service>.exe` for Windows.
- **air-%**: Runs **Air** for a specific service (e.g., `make air-api-gateway` or `make air-bus-service`) in its directory.
- **air-all**: Starts **Air** for all services in parallel (use with caution due to resource usage).
- **docker**: Starts Docker Compose to run MySQL, Redis, Kafka, etc.
- **docker-down**: Stops Docker Compose.
- **test**: Runs tests in the `tests/` directory.
- **clean**: Removes compiled binaries and log files.
- **deps**: Installs Go dependencies and **Air**.

### Prerequisites
Ensure the following are in place:
- **Air** is installed: `go install github.com/air-verse/air@latest`.
- `$GOPATH/bin` (e.g., `C:\Users\ACER\go\bin`) is in your PATH:
  ```bash
  export PATH=$PATH:$GOPATH/bin
  echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
  source ~/.bashrc
  ```
- **Air** configs exist for each service (e.g., `cmd/api-gateway/.air.toml`, `cmd/bus-service/.air.toml`) with `.exe` for Windows, as provided earlier:
  ```toml
  root = "."
  tmp_dir = "../../build/tmp"

  [build]
    bin = "../../build/bin/api-gateway.exe"
    cmd = "go build -o ../../build/bin/api-gateway.exe ./main.go"
    include_ext = ["go"]
    exclude_dir = ["../../build", "../../migrations", "../../tests"]
    delay = 1000

  [log]
    time = true
  ```
- Logger is set up (`pkg/logger/logger.go`) to write to `build/tmp/logs/<service-name>.log`, as provided previously:
  ```go
  func Init(serviceName string) *Logger {
      logDir := "../../build/tmp/logs"
      if err := os.MkdirAll(logDir, 0755); err != nil {
          log.Fatalf("Failed to create log directory: %v", err)
      }
      logFilePath := filepath.Join(logDir, serviceName+".log")
      file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
      if err != nil {
          log.Fatalf("Failed to open log file: %v", err)
      }
      writer := bufio.NewWriter(file)
      multiWriter := io.MultiWriter(os.Stdout, writer)
      logger := log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)
      return &Logger{Logger: logger, writer: writer, file: file}
  }
  ```

### Usage
Run these commands from `D:\GIN\real-time-bus-tracking-service`:

1. **Build all services**:
   ```bash
   make build
   ```
   Outputs: `build/bin/api-gateway.exe`, `build/bin/bus-service.exe`, etc.

2. **Run API Gateway with Air**:
   ```bash
   make air-api-gateway
   ```
   Starts **Air** in `cmd/api-gateway/`, hot reloading on `.go` changes, with logs in `build/tmp/logs/api-gateway.log`.

3. **Run Bus Service with Air**:
   ```bash
   make air-bus-service
   ```
   Starts **Air** in `cmd/bus-service/`, with logs in `build/tmp/logs/bus-service.log`.

4. **Run all services with Air** (in separate terminals or background):
   ```bash
   make air-all
   ```

5. **Start Docker dependencies**:
   ```bash
   make docker
   ```

6. **Stop Docker**:
   ```bash
   make docker-down
   ```

7. **Run tests**:
   ```bash
   make test
   ```

8. **Clean artifacts**:
   ```bash
   make clean
   ```

9. **Install dependencies**:
   ```bash
   make deps
   ```

### Test Logs
Run the **API Gateway**:
```bash
make air-api-gateway
```

Trigger a log entry:
```bash
curl -X POST http://localhost:8080/buses -d '{"license_plate":"ABC123","wheelchair_enabled":true}'
```

Check the log file:
```bash
cat D:\GIN\real-time-bus-tracking-service\build\tmp\logs\api-gateway.log
```

**Expected Output** (line-by-line, as ensured by `bufio.Writer`):
```
2025/07/29 19:30:00 handlers/bus_handlers.go:20: [INFO] Bus created [license_plate ABC123]
```

### Notes
- **Windows Paths**: The **Makefile** uses `D:/GIN/real-time-bus-tracking-service` for Windows compatibility in MinGW64. If your path differs, update `ROOT_DIR`.
- **Service List**: The `SERVICES` variable includes all 12 microservices. Remove any that don’t exist yet (e.g., `route-service`) or add new ones as needed.
- **Log Files**: Each service’s logs are written to `build/tmp/logs/<service-name>.log` due to `logger.Init("<service-name>")`.
- **Air Configs**: Ensure each service has a `.air.toml` with `.exe` (e.g., `cmd/bus-service/.air.toml` from previous responses).
- **Docker**: Run `make docker` before `make air-<service>` if services depend on MySQL, Redis, or Kafka.

If you need specific **Makefile** targets for other services, additional tasks (e.g., linting), or help with a specific microservice, let me know!