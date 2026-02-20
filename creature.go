package main

import (
	"math/rand"
)

type Vec2 struct {
	X, Y int
}

type Creature struct {
	ID  int
	Pos Vec2
}

func (c *Creature) Update(w *World) {
	dx := rand.Intn(3) - 1
	dy := rand.Intn(3) - 1

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
