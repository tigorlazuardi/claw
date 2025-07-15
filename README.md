# Claw

A cross-platform image downloader and collection manager designed to fetch images from various sources (Reddit, Booru sites, etc.) and organize them for easy syncing to multiple devices.

## Overview

Claw is a client-server application that automatically downloads and manages images from various online sources. It's designed to collect wallpapers and images based on user-defined criteria and organize them in a way that facilitates easy syncing to different devices using tools like Syncthing.

## Features

- **Multi-source Support**: Download images from Reddit, Booru sites, and other configurable sources
- **Device-aware Collection**: Filter and organize images based on device constraints (resolution, aspect ratio, file size)
- **Scheduled Downloads**: Automatic image collection using cron-based scheduling
- **Cross-platform Clients**: Web interface, mobile apps, and CLI daemon
- **Smart Deduplication**: Avoid re-downloading existing images
- **Sync-friendly Storage**: Uses hard links when possible to optimize storage

## Architecture

### Server (Single Binary)
- **Language**: Go 1.24.4
- **Database**: SQLite (pure Go implementation)
- **Migrations**: Goose
- **Query Builder**: go-jet/jet
- **Communication**: ConnectRPC
- **Constraints**: Pure Go libraries only (no CGO)

### Clients
- **Web UI**: Svelte (image viewing, device/source management)
- **Mobile**: CapacitorJS (wallpaper selection and scheduling)
- **CLI**: Daemon mode (event hooks, wallpaper management)

## Core Concepts

### Source
Defines how to obtain images from a specific platform or service:
- Configurable parameters (subreddits, usernames, etc.)
- Countback parameter to limit search depth
- Handles authentication and rate limiting

### Device
Defines constraints for image selection:
- Resolution requirements (min/max width/height)
- Aspect ratio tolerance
- File size limits
- NSFW content handling

### Image
Represents a downloaded image with metadata:
- Dimensions and file size
- Source information
- Device assignments
- User tags and favorites

### Worker
Manages scheduled image downloads:
- Cron-based scheduling
- Executes source parameters
- Handles download failures and retries

## Development

### Prerequisites
- Go 1.24.4+
- Nix (recommended for development environment)

### Setup
```bash
# Clone the repository
git clone https://github.com/tigorlazuardi/claw.git
cd claw

# Enter development environment (with Nix)
nix develop

# Initialize dependencies
go mod tidy

# Run database migrations
go run ./cmd/goose
```

### Development Commands
```bash
# Build the project
go build ./...

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run

# Run migrations
go run ./cmd/goose
```

## Project Structure

```
├── cmd/
│   ├── claw/          # Main CLI application
│   ├── goose/         # Database migration tool
│   └── jet/           # Code generation for go-jet
├── lib/claw/          # Core API implementations
│   └── types/         # Custom types and serialization
├── server/            # Server code (routers, handlers)
├── worker/            # Background worker implementations
├── ui/                # Svelte UI and CapacitorJS mobile apps
├── migrations/        # Database migration files
└── artifacts/         # Development artifacts (database, etc.)
```

## Configuration

The application uses environment variables and configuration files:

- `GOOSE_DBSTRING`: Database connection string for migrations
- Additional configuration will be documented as the project develops

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

Please follow conventional commit format with scope (e.g., `feat(server): add image download endpoint`).

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Status

🚧 **Early Development** - This project is in active development. Features and APIs may change.

## Roadmap

- [ ] Core server implementation
- [ ] Basic image source adapters (Reddit, Booru)
- [ ] Device management system
- [ ] Web UI implementation
- [ ] Mobile app development
- [ ] CLI daemon functionality
- [ ] Advanced scheduling and filtering
- [ ] Documentation and examples