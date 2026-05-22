# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A life simulation server in Go that simulates creatures moving in a 2D world. Creatures randomly move, consume energy, and die when energy depletes. The server exposes the world state via HTTP REST and WebSocket APIs.

## Commands

**Run the server:**
```bash
go run .
```
Server starts at http://localhost:8080

**Build:**
```bash
go build
```

**Install dependencies:**
```bash
go mod download
```

## Architecture

### Core Components

- **World** (`main.go`): Central simulation state containing all creatures. Runs at 20 ticks/second in a goroutine with mutex-protected concurrent access (RWMutex for read operations, full Lock for writes).

- **Creature** (`creature.go`): Entities that move randomly each tick, gain random energy (food), and lose 0.5 energy per tick. Die when energy ≤ 0.

- **Server** (`server.go`): HTTP handlers for querying world state
  - `/region` - HTTP GET endpoint with query params `x1`, `y1`, `x2`, `y2` returns creatures in that region as JSON
  - `/ws` - WebSocket endpoint accepts JSON messages with type `get_region` and payload with region coordinates

### Concurrency Model

The World runs independently in a goroutine (`world.Run()`). HTTP handlers query the world using read locks (`RLock`) for `GetRegion`. The world simulation uses full locks (`Lock`) during `Step()` when updating all creatures.

### Configuration Constants

In `main.go`:
- `TickRate = 20` - World updates per second
- `CPUProfiler = false` - CPU profiling toggle (not currently implemented)
- World dimensions: 1000x1000
- Initial creatures: 10,000

## Dependencies

- `github.com/gorilla/websocket` - WebSocket support with permissive CORS for development

## Development Rules

- **Language**: Use English for ALL code, comments, documentation, commit messages, and files in this project
- **Git workflow**: Commit ALL changes with descriptive commit messages. Never leave uncommitted work
