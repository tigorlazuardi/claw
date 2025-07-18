# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project named "claw" (github.com/tigorlazuardi/claw) using Go version 1.24.4.

## Current State

The project appears to be in its initial state with only a `go.mod` file present. No source code, documentation, or build configuration has been created yet.

## Development Commands

Since this is a fresh Go project, standard Go commands will be used:

```bash
# Initialize and download dependencies
go mod tidy

# Build the project
go build ./...

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestFunctionName ./...

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Run the main application (once main.go is created)
go run main.go

# Run database migrations
go run ./cmd/goose
```

## Build Commands

- You may use `go build -o artifacts/claw ./cmd/claw/main.go` to test builds

## Architecture Notes

- The project is a Client-Server architecture
- The server aims for a single binary for easy distribution
- The client must use cross platform technologies for less work being done in the UI side

## Project Purpose

- This project aims to create a downloader and collect images from various sources, like Reddit and Booru sites
- Takes into account the user's target devices
- Collects and groups images into a directory to facilitate easy syncing to devices (e.g., using Syncthing)

## Technical Constraints and Library Choices

- Because of the goal of single binary server, the project must use pure go libs (don't use CGO)
- The primary backend database is sqlite, and must use pure go sqlite database
- goose library will be used for migrations and go-jet/jet library will be used as query builder to answer dynamic queries from users
- ConnectRPC will be used as primary communication protocol between client and server
- Since this projects aims to support browsers, Services must not use BiDi streamings, only Unary and single side streaming is allowed

## Client Implementation Details

### Framework and Technologies

- Primary client framework: Svelte (NOT SvelteKit)
- Mobile deployment: CapacitorJS
- Supports both web and mobile interfaces

### Web UI Capabilities

- Simple viewer for collected images
- Allows moderation of devices and sources

### Mobile Device Capabilities

- Extends web UI functionality
- Enables user to:
  - Select collected wallpapers
  - Change device wallpaper
  - Schedule wallpaper changes

### CLI Client Interface

- Runs as a daemon
- Supports:
  - Receiving event hooks
  - Listing wallpapers
  - User-specified commands for wallpaper management

### Platform Differences

- Web UI: Basic image viewing and device/source management
- Mobile Devices: Advanced wallpaper selection and scheduling features
- CLI: Background daemon with event and wallpaper management capabilities

## Core Project Concepts

### Source

- Source determines how to obtain images and handles image retrieval
- Receives parameters to target specific images (e.g., subreddits or usernames)
- Parameters can be simple strings or complex like JSON strings
- Responsible for interpreting and processing input parameters
- Includes `countback` parameter to specify how far back to look for images
- Counts posts even if they are not images to avoid infinite searching

### Device

- Defines constraints for image selection and device assignment
- Contains filters like:
  - Height and width constraints
  - Aspect ratio tolerance
  - Minimum and maximum height/width
  - NSFW image handling (accept, block, or only NSFW)
  - Minimum and maximum file size
- Aspect ratio tolerance calculates image "shape" match by comparing height/width ratios

### Image

- Represents an image with detailed metadata
- Contains:
  - Height and width
  - File size
  - Source origin
  - Download location
  - Title and description (if available)
  - Assigned devices
- Supports user tagging and favoriting
- Tracks file copies across filesystem
- Ensures sync-friendly operation by:
  - Attempting hard links first when copying to multiple devices
  - Avoiding re-downloading existing images on subsequent runs

### Worker

- Manages scheduled image downloads
- Sources receive a `schedule` parameter with cron expressions
- When a schedule matches, a Worker is assigned
- Worker initiates image download using existing source parameters

## Project Structure and Key Components

### Project Directory Structure

- `cmd/claw` contains main entry for cli code. The library will use urfave/cli/v3 to handle command line parsing.
- `migrations` dir contains migration sql files, and one `migration.go` that will execute migration files that will be embedded using embed package.
- `cmd/goose` calls the migrations dir so the developer can migrate on demand without running the application.
- `cmd/jet` runs customized code generation for go-jet, because we will use custom types for more sane interaction between code and database. e.g. if a table field contains `is_` prefix and type of `INTEGER`, it will use `types.Bool` as bridge between database 0/1 value and treated as proper boolean in code side.
- `lib/claw/types` contains custom types to streamlines how a code serializes and deserializes. e.g. `types.Bool` above.

## Folder Descriptions

- `ui` folder contains the user interface codes (the Svelte code) and CapacitorJS to create mobile apps.

## Project Directory Details

- `lib/claw` contains implementations of APIs.
- `server` contains server codes (e.g. Routers and Handlers).
- `worker` contains worker codes.

## Library References

- library documentation for `urfave/cli/v3` is located at https://cli.urfave.org/v3/getting-started. Refer to use them when creating cli commands.
- documentation for mapping flags to env for urfave/cli/v3 is at https://cli.urfave.org/v3/examples/flags/value-sources/

## Best Practices and Guidelines

- CLAUDE MUST NOT self promote in git commits
- Git commit style should use conventional commit with scope
- When creating git commits, you are allowed to run multiple git adds and git commits so the changes are atomic.
- When creating new structs, types, or functions, always include documentation for them

## Database Field Naming Conventions

- When creating fields for database and migrations use the following rules:
  1. Timestamp related fields must have suffix `_at` and uses INTEGER value. The actual value is Unix milliseconds from epoch.
  2. Duration related fields must have suffix `_dur` and uses INTEGER value. The maximum precision value is milliseconds, so the database value is always using milliseconds going in and out.
  3. For boolean kind of field, field must have `is_` prefix and uses INTEGER value.

## Dependencies

- When dealing with protobuf dependencies, assume all dependencies are available locally.

## Protobuf Conventions

- When dealing with protobuf schemas, follow the conventions that buf.build recommends. Run `buf lint` to ensure no warnings or errors occurs.
- When dealing with probouf schemas, you are allowed to run `buf generate` on `schemas` folder to validate the output or whenever it feels it's required to do so.
- When creating protos, to avoid CRC errors, make sure the proto files are under `schemas/claw/v1` and the packages have prefix of `claw.v1` to follow the buf build guideline. All proto files dedicated to this project is under that one folder, so make sure to use proper name prefixes or suffixes for messages, services, etc to keep unique.