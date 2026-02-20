package main

import (
	"math/rand"
)

type Vec2 struct {
	X, Y int
}

type Creature struct {
	ID     int
	Pos    Vec2
	Energy float32
}

func (c *Creature) Update(w *World) {
	dx := rand.Intn(3) - 1
	dy := rand.Intn(3) - 1
	food := rand.Float32()

	c.Energy += food
	c.Energy -= 0.5

	if c.Energy <= 0 {
		w.RemoveCreature(c.ID)
		return
	}

	c.Pos.X += dx
	c.Pos.Y += dy

	if c.Pos.X < 0 {
		c.Pos.X = 0
	}
	if c.Pos.Y < 0 {
		c.Pos.Y = 0
	}
	if c.Pos.X >= w.Width {
		c.Pos.X = w.Width
	}
	if c.Pos.Y >= w.Height {
		c.Pos.Y = w.Height
	}
}
