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
	WorldWidth           int     `json:"world_width"`
	WorldHeight          int     `json:"world_height"`
	InitialCreatures     int     `json:"initial_creatures"`
	CreatureInitialEnergy float32 `json:"creature_initial_energy"`
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

	Tick      int
	Creatures map[int]*Creature
	Width     int
	Height    int
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

func (w *World) RemoveCreature(id int) {
	delete(w.Creatures, id)
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

	world := &World{
		Width:     config.WorldWidth,
		Height:    config.WorldHeight,
		Creatures: make(map[int]*Creature),
	}
	for i := range config.InitialCreatures {
		world.Creatures[i] = &Creature{
			ID:     i,
			Pos:    Vec2{rand.Intn(world.Width), rand.Intn(world.Height)},
			Energy: config.CreatureInitialEnergy,
		}
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
