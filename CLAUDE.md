# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A life simulation server in Go that simulates creatures moving in a 2D world with procedurally generated terrain. Creatures randomly move on plains, consume energy, die when energy depletes, and can reproduce when well-fed. The simulation includes overpopulation mechanics that affect food availability. Terrain is generated using Perlin noise and includes water, plains, and mountains. The server exposes the world state via HTTP REST and WebSocket APIs.

## Configuration

Simulation parameters are loaded from `config.json` at startup:
- `world_width`, `world_height` - World dimensions
- `initial_creatures` - Number of creatures to spawn
- `creature_initial_energy` - Starting energy for each creature
- `max_population` - Maximum sustainable population (affects food multiplier)
- `reproduction_energy_threshold` - Energy required to reproduce
- `reproduction_energy_cost` - Energy cost for parent when reproducing
- `terrain_seed` - Seed for procedural terrain generation (deterministic)

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

- **World** (`main.go`): Central simulation state containing all creatures and terrain generator. Runs at 20 ticks/second in a goroutine with mutex-protected concurrent access (RWMutex for read operations, full Lock for writes).

- **Terrain** (`terrain.go`): Procedural terrain generator using Perlin noise. Generates three terrain types (Water=0, Plains=1, Mountains=2) deterministically from seed. Water borders (50px) around map edges. Full map is pre-generated at startup and stored in World.TerrainMap for fast lookups.

- **Creature** (`creature.go`): Entities that move randomly on plains each tick, gain energy (food) based on population density with complex distribution (exponential penalty for overpopulation, high food chance at low population), and lose 0.5 energy per tick. Cannot move into water/mountains. Die when energy ≤ 0. Reproduce when energy >= threshold by spawning offspring on nearby walkable terrain.

- **Server** (`server.go`): HTTP handlers for querying world state
  - `/region` - HTTP GET endpoint with query params `x1`, `y1`, `x2`, `y2` returns creatures and terrain data in that region as JSON. Includes CORS headers and coordinate validation.
  - `/stats` - HTTP GET endpoint that returns population statistics including total creatures and energy distribution (low/medium/high)
  - `/ws` - WebSocket endpoint that sends `world_bounds` message on connection, then accepts JSON messages:
    - `get_region` - returns creatures and terrain in specified region
    - `get_stats` - returns population statistics
    - `get_map` - returns full terrain map (entire world)

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
