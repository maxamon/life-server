# Life Simulation Server

Go server that simulates creatures moving in a 2D world. Creatures randomly move, consume energy, and die when energy depletes.

## Requirements

- Go 1.22.5+

## Configuration

Edit `config.json` to adjust simulation parameters:

```json
{
  "world_width": 1000,
  "world_height": 1000,
  "initial_creatures": 10000,
  "creature_initial_energy": 10.0
}
```

## Quick Start

```bash
# Install dependencies
go mod download

# Run server
go run .
```

Server starts at http://localhost:8080

## Build

```bash
go build
./life-server
```

## API

### HTTP REST

**GET /region** - Get creatures in region

Query params:
- `x1`, `y1` - top-left corner
- `x2`, `y2` - bottom-right corner

Example:
```bash
curl "http://localhost:8080/region?x1=0&y1=0&x2=100&y2=100"
```

### WebSocket

**WS /ws** - WebSocket connection for real-time data

On connection, server sends world bounds:
```json
{
  "type": "world_bounds",
  "payload": {
    "min_x": 0,
    "min_y": 0,
    "max_x": 1000,
    "max_y": 1000
  }
}
```

Client message format:
```json
{
  "type": "get_region",
  "payload": {
    "x1": 0,
    "y1": 0,
    "x2": 100,
    "y2": 100
  }
}
```

## Simulation Parameters

Default values (configurable in `config.json`):
- World size: 1000x1000
- Initial creatures: 10,000
- Starting energy: 10.0

Fixed parameters:
- Update rate: 20 ticks/sec
- Energy loss: 0.5/tick

## Architecture

- **World** - central simulation state with concurrent access (RWMutex)
- **Creature** - entities with position and energy
- **Server** - HTTP/WebSocket handlers

Simulation runs in separate goroutine with mutex synchronization and graceful shutdown support.
