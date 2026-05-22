package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	TickRate = 20
)

type Config struct {
	WorldWidth                  int     `json:"world_width"`
	WorldHeight                 int     `json:"world_height"`
	InitialCreatures            int     `json:"initial_creatures"`
	CreatureInitialEnergy       float32 `json:"creature_initial_energy"`
	MaxPopulation               int     `json:"max_population"`
	ReproductionEnergyThreshold float32 `json:"reproduction_energy_threshold"`
	ReproductionEnergyCost      float32 `json:"reproduction_energy_cost"`
	TerrainSeed                 int64   `json:"terrain_seed"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

type World struct {
	sync.RWMutex

	Tick                        int
	Creatures                   map[int]*Creature
	Width                       int
	Height                      int
	NextCreatureID              int
	MaxPopulation               int
	ReproductionEnergyThreshold float32
	ReproductionEnergyCost      float32
	CreatureInitialEnergy       float32
	Terrain                     *TerrainGenerator
}

func (w *World) Step() {
	w.Lock()
	defer w.Unlock()

	w.Tick++

	for _, c := range w.Creatures {
		c.Update(w)
	}
}

func (w *World) GetRegion(x1, y1, x2, y2 int) []*Creature {
	w.RLock()
	defer w.RUnlock()

	var result []*Creature

	for _, c := range w.Creatures {
		if c.Pos.X >= x1 && c.Pos.X <= x2 &&
			c.Pos.Y >= y1 && c.Pos.Y <= y2 {
			result = append(result, c)
		}
	}
	return result
}

func (w *World) GetRegionWithTerrain(x1, y1, x2, y2 int) (creatures []*Creature, terrain [][]byte) {
	w.RLock()
	defer w.RUnlock()

	// Get creatures
	for _, c := range w.Creatures {
		if c.Pos.X >= x1 && c.Pos.X <= x2 &&
			c.Pos.Y >= y1 && c.Pos.Y <= y2 {
			creatures = append(creatures, c)
		}
	}

	// Get terrain
	height := y2 - y1 + 1
	width := x2 - x1 + 1
	terrain = make([][]byte, height)
	for dy := 0; dy < height; dy++ {
		terrain[dy] = make([]byte, width)
		for dx := 0; dx < width; dx++ {
			terrain[dy][dx] = byte(w.Terrain.GetTerrain(x1+dx, y1+dy))
		}
	}

	return creatures, terrain
}

func (w *World) RemoveCreature(id int) {
	delete(w.Creatures, id)
}

func (w *World) AddCreature(pos Vec2, energy float32) {
	id := w.NextCreatureID
	w.NextCreatureID++
	w.Creatures[id] = &Creature{
		ID:     id,
		Pos:    pos,
		Energy: energy,
	}
}

func (w *World) GetStats() map[string]interface{} {
	w.RLock()
	defer w.RUnlock()

	total := len(w.Creatures)
	lowEnergy := 0    // < 3
	mediumEnergy := 0 // >= 3 and < 7
	highEnergy := 0   // >= initial energy

	for _, c := range w.Creatures {
		if c.Energy < 3 {
			lowEnergy++
		} else if c.Energy < 7 {
			mediumEnergy++
		}
		if c.Energy >= w.CreatureInitialEnergy {
			highEnergy++
		}
	}

	return map[string]interface{}{
		"tick":                  w.Tick,
		"total_creatures":       total,
		"low_energy_count":      lowEnergy,
		"medium_energy_count":   mediumEnergy,
		"high_energy_count":     highEnergy,
		"max_population":        w.MaxPopulation,
		"population_percentage": float32(total) / float32(w.MaxPopulation) * 100,
	}
}

func (w *World) Run(ctx context.Context) {
	timer := time.NewTicker(time.Second / TickRate)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			w.Step()
		}
	}
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	terrain := NewTerrainGenerator(config.TerrainSeed, config.WorldWidth, config.WorldHeight)

	world := &World{
		Width:                       config.WorldWidth,
		Height:                      config.WorldHeight,
		Creatures:                   make(map[int]*Creature),
		NextCreatureID:              config.InitialCreatures,
		MaxPopulation:               config.MaxPopulation,
		ReproductionEnergyThreshold: config.ReproductionEnergyThreshold,
		ReproductionEnergyCost:      config.ReproductionEnergyCost,
		CreatureInitialEnergy:       config.CreatureInitialEnergy,
		Terrain:                     terrain,
	}

	// Spawn creatures only on plains
	spawned := 0
	maxAttempts := config.InitialCreatures * 10
	for attempts := 0; spawned < config.InitialCreatures && attempts < maxAttempts; attempts++ {
		x := rand.Intn(world.Width)
		y := rand.Intn(world.Height)

		if terrain.IsWalkable(x, y) {
			world.Creatures[spawned] = &Creature{
				ID:     spawned,
				Pos:    Vec2{x, y},
				Energy: config.CreatureInitialEnergy,
			}
			spawned++
		}
	}

	if spawned < config.InitialCreatures {
		log.Printf("Warning: only spawned %d/%d creatures", spawned, config.InitialCreatures)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go world.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.Write([]byte(`<html><body><h1>Hello world</h1></body></html>`))
	})
	mux.HandleFunc("/region", regionHandler(world))
	mux.HandleFunc("/stats", statsHandler(world))
	mux.HandleFunc("/ws", wsHandler(world))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Started server http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
