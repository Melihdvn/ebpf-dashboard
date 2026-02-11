# eBPF Dashboard Backend

Real-time application performance monitoring backend using Go and eBPF/BCC tools.

## Features

- **Process Monitoring**: Track process execution events using `execsnoop`
- **Network Monitoring**: Monitor TCP connections using `tcpconnect`
- **Disk I/O Monitoring**: Analyze disk latency distribution using `biolatency`
- **CPU Profiling**: Collect CPU stack traces for flame graph visualization using `profile-bpfcc`
- **REST API**: Clean RESTful API for accessing metrics
- **SQLite Storage**: Persistent storage of all collected metrics
- **Real-time Collection**: Background collectors running continuously

## Architecture

```
ebpf-dashboard/
├── database/          # Database initialization
├── models/            # Data models
├── collector/         # BCC tool collectors
├── repository/        # Data access layer
├── services/          # Business logic
├── handlers/          # HTTP API handlers
├── utils/             # Shared utilities
└── main.go            # Application entry point
```

## Prerequisites

- Go 1.24+
- BCC tools installed:
  - `execsnoop`
  - `tcpconnect`
  - `biolatency`
- Sudo privileges (required for eBPF)

## Installation

```bash
# Clone or navigate to the project
cd /home/melih/Desktop/ebpf-dashboard

# Install dependencies
go mod download
```

## Running the Application

**Important**: The application must be run with sudo privileges because BCC tools require root access.

```bash
sudo go run main.go
```

You should see:
```
Starting eBPF Dashboard Backend...
Database initialized successfully
Process collector started
Network collector started
Disk latency collector started
Server is running on http://localhost:8080
```

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Get Process Events
```bash
# Get last 50 processes (default)
curl http://localhost:8080/api/metrics/processes

# Get last 10 processes
curl http://localhost:8080/api/metrics/processes?limit=10
```

### Get Network Connections
```bash
# Get last 50 connections (default)
curl http://localhost:8080/api/metrics/network

# Get last 20 connections
curl http://localhost:8080/api/metrics/network?limit=20
```

### Get Disk I/O Latency
```bash
# Get latest latency histogram
curl http://localhost:8080/api/metrics/disk

# Get last 10 histogram buckets
curl http://localhost:8080/api/metrics/disk?limit=10
```

### Get CPU Profiling Data
```bash
# Get last 50 CPU profile samples (default)
curl http://localhost:8080/api/metrics/cpuprofile

# Get last 20 samples
curl http://localhost:8080/api/metrics/cpuprofile?limit=20
```

## Data Collection

The application runs four background collectors:

- **Process Collector**: Runs `execsnoop` continuously, streams events in real-time
- **Network Collector**: Runs `tcpconnect` continuously, captures TCP connections as they happen
- **Disk Collector**: Runs `biolatency` every 5 seconds to collect I/O latency histograms
- **CPU Profile Collector**: Runs `profile-bpfcc` every 5 seconds to collect CPU stack traces for flame graph visualization

Process and network events are captured immediately as they occur and saved to the database every second. This provides true real-time monitoring of system activity.

All data is stored in `metrics.db` (SQLite database).

## Graceful Shutdown

Press `Ctrl+C` to stop the server. The application will:
1. Stop all background collectors
2. Wait for goroutines to finish
3. Close database connections
4. Exit cleanly

## Project Structure

- **database/**: Database initialization and schema management
- **models/**: Data structures for Process, Network, and Disk metrics
- **collector/**: BCC tool integration and output parsing
- **repository/**: Database operations (CRUD)
- **services/**: Business logic and background collection
- **handlers/**: HTTP request handlers
- **utils/**: Sudo execution helper

## Development

### Building
```bash
go build -o ebpf-dashboard
sudo ./ebpf-dashboard
```

### Adding New Metrics

To add a new metric type, follow the pattern:

1. Create model in `models/`
2. Create collector in `collector/`
3. Create repository in `repository/`
4. Create service in `services/`
5. Create handler in `handlers/`
6. Register routes in `main.go`

## License

This project is for educational purposes as part of an internship project.
