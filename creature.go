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
	// Movement
	dx := rand.Intn(3) - 1
	dy := rand.Intn(3) - 1

	// Food calculation based on population
	currentPopulation := len(w.Creatures)
	foodMultiplier := float32(1.0)

	if currentPopulation > w.MaxPopulation {
		// Overpopulation - less food
		foodMultiplier = float32(w.MaxPopulation) / float32(currentPopulation)
	} else if currentPopulation < w.MaxPopulation {
		// Underpopulation - more food
		foodMultiplier = 1.0 + (float32(w.MaxPopulation-currentPopulation) / float32(w.MaxPopulation) * 0.5)
	}

	food := rand.Float32() * foodMultiplier

	c.Energy += food
	c.Energy -= 0.5

	// Death
	if c.Energy <= 0 {
		w.RemoveCreature(c.ID)
		return
	}

	// Reproduction
	if c.Energy >= w.ReproductionEnergyThreshold {
		// Find position near parent
		newPos := Vec2{
			X: c.Pos.X + rand.Intn(3) - 1,
			Y: c.Pos.Y + rand.Intn(3) - 1,
		}

		// Clamp to world bounds
		if newPos.X < 0 {
			newPos.X = 0
		}
		if newPos.Y < 0 {
			newPos.Y = 0
		}
		if newPos.X >= w.Width {
			newPos.X = w.Width - 1
		}
		if newPos.Y >= w.Height {
			newPos.Y = w.Height - 1
		}

		// Create new creature
		w.AddCreature(newPos, w.CreatureInitialEnergy)

		// Parent loses energy
		c.Energy -= w.ReproductionEnergyCost
	}

	// Move
	c.Pos.X += dx
	c.Pos.Y += dy

	// Boundary check
	if c.Pos.X < 0 {
		c.Pos.X = 0
	}
	if c.Pos.Y < 0 {
		c.Pos.Y = 0
	}
	if c.Pos.X >= w.Width {
		c.Pos.X = w.Width - 1
	}
	if c.Pos.Y >= w.Height {
		c.Pos.Y = w.Height - 1
	}
}
