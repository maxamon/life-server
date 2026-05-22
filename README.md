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
  "creature_initial_energy": 10.0,
  "max_population": 15000,
  "reproduction_energy_threshold": 15.0,
  "reproduction_energy_cost": 8.0
}
```

**Parameters:**
- `world_width`, `world_height` - World dimensions
- `initial_creatures` - Starting number of creatures
- `creature_initial_energy` - Energy for new creatures
- `max_population` - Maximum sustainable population (affects food availability)
- `reproduction_energy_threshold` - Energy needed to reproduce
- `reproduction_energy_cost` - Energy parent loses when reproducing

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

**GET /stats** - Get population statistics

Returns:
```json
{
  "tick": 1234,
  "total_creatures": 9876,
  "low_energy_count": 123,
  "medium_energy_count": 456,
  "high_energy_count": 789,
  "max_population": 15000,
  "population_percentage": 65.84
}
```

- `low_energy_count` - creatures with energy < 3
- `medium_energy_count` - creatures with 3 <= energy < 7
- `high_energy_count` - creatures with energy >= initial energy

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

Client message formats:

Get creatures in region:
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

Get population statistics:
```json
{
  "type": "get_stats"
}
```

Server responds with:
```json
{
  "type": "stats",
  "payload": {
    "tick": 1234,
    "total_creatures": 9876,
    ...
  }
}
```

## Simulation Mechanics

**Energy System:**
- Creatures lose 0.5 energy per tick
- Food gain varies based on population:
  - Overpopulation (> max): reduced food
  - Underpopulation (< max): increased food
- Death occurs when energy <= 0

**Reproduction:**
- Creatures with energy >= threshold can reproduce
- Offspring spawns at nearby position
- Parent loses energy cost for reproduction
- Child starts with initial energy

**Fixed Parameters:**
- Update rate: 20 ticks/sec
- Energy loss: 0.5/tick

## Architecture

- **World** - central simulation state with concurrent access (RWMutex)
- **Creature** - entities with position and energy
- **Server** - HTTP/WebSocket handlers

Simulation runs in separate goroutine with mutex synchronization and graceful shutdown support.
