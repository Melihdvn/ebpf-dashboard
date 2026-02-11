# eBPF Dashboard

Real-time application performance monitoring system using eBPF technology.

## Project Structure

```
ebpf-dashboard/
├── backend/          # Go backend with eBPF collectors
└── frontend/         # (Coming soon) Web dashboard UI
```

## Backend

The backend is built with Go and uses BCC (BPF Compiler Collection) tools to collect system metrics in real-time.

**Features:**
- Process execution monitoring (execsnoop)
- TCP connection tracking (tcpconnect)
- Disk I/O latency analysis (biolatency)
- REST API for data access
- SQLite database for persistence

See [backend/README.md](backend/README.md) for detailed documentation.

## Quick Start

### Prerequisites
- Go 1.24+
- BCC tools installed
- Sudo privileges

### Running the Backend

```bash
cd backend
sudo go run main.go
```

The API will be available at `http://localhost:8080`

### API Endpoints

- `GET /health` - Health check
- `GET /api/metrics/processes` - Process execution events
- `GET /api/metrics/network` - TCP connections
- `GET /api/metrics/disk` - Disk I/O latency

## Development

This project is part of an internship focused on learning eBPF technology and building real-world monitoring applications.

## License

Educational project for internship purposes.
