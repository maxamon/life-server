# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A life simulation server in Go that simulates creatures moving in a 2D world. Creatures randomly move, consume energy, die when energy depletes, and can reproduce when well-fed. The simulation includes overpopulation mechanics that affect food availability. The server exposes the world state via HTTP REST and WebSocket APIs.

## Configuration

Simulation parameters are loaded from `config.json` at startup:
- `world_width`, `world_height` - World dimensions
- `initial_creatures` - Number of creatures to spawn
- `creature_initial_energy` - Starting energy for each creature
- `max_population` - Maximum sustainable population (affects food multiplier)
- `reproduction_energy_threshold` - Energy required to reproduce
- `reproduction_energy_cost` - Energy cost for parent when reproducing

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

- **Creature** (`creature.go`): Entities that move randomly each tick, gain random energy (food) affected by population density, and lose 0.5 energy per tick. Die when energy ≤ 0. Reproduce when energy >= threshold by spawning offspring nearby.

- **Server** (`server.go`): HTTP handlers for querying world state
  - `/region` - HTTP GET endpoint with query params `x1`, `y1`, `x2`, `y2` returns creatures in that region as JSON. Includes CORS headers and coordinate validation.
  - `/stats` - HTTP GET endpoint that returns population statistics including total creatures and energy distribution (low/medium/high)
  - `/ws` - WebSocket endpoint that sends `world_bounds` message on connection, then accepts JSON messages:
    - `get_region` - returns creatures in specified region
    - `get_stats` - returns population statistics

### Concurrency Model

The World runs independently in a goroutine (`world.Run()`). HTTP handlers query the world using read locks (`RLock`) for `GetRegion`. The world simulation uses full locks (`Lock`) during `Step()` when updating all creatures.

### Configuration

In `main.go`:
- `TickRate = 20` - World updates per second (constant)

In `config.json` (loaded at startup):
- `world_width`, `world_height` - World dimensions (default: 1000x1000)
- `initial_creatures` - Number of creatures to spawn (default: 10,000)
- `creature_initial_energy` - Starting energy (default: 10.0)

## Dependencies

- `github.com/gorilla/websocket` - WebSocket support with permissive CORS for development

## Development Rules

- **Language**: Use English for ALL code, comments, documentation, commit messages, and files in this project
- **Git workflow**: Commit ALL changes with descriptive commit messages. Never leave uncommitted work
