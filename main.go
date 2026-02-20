package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	TickRate    = 20
	CPUProfiler = false
)

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

func (w *World) Run() {
	timer := time.NewTicker(time.Second / TickRate)
	defer timer.Stop()

	for range timer.C {
		w.Step()
	}
}

func log(data string) {
	fmt.Println(data)
}

func main() {
	world := &World{
		Width:     1000,
		Height:    1000,
		Creatures: make(map[int]*Creature),
	}
	for i := range 10000 {
		world.Creatures[i] = &Creature{
			ID:     i,
			Pos:    Vec2{i % 1000, i % 1000},
			Energy: 10.0,
		}
	}

	go world.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.Write([]byte(`<html><body><h1>Hello world</h1></body></html>`))
	})
	mux.HandleFunc("/region", regionHandler(world))
	mux.HandleFunc("/ws", wsHandler(world))
	fmt.Println("Started server http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
