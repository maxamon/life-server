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
	dx := 2 * (rand.Intn(3) - 1)
	dy := 2 * (rand.Intn(3) - 1)

	// Food calculation based on population
	currentPopulation := len(w.Creatures)
	var food float32

	populationRatio := float32(currentPopulation) / float32(w.MaxPopulation)

	if currentPopulation > w.MaxPopulation {
		// Overpopulation - exponential food decrease
		// The more overpopulated, the less food
		overpopulationFactor := populationRatio - 1.0 // 0.0 to infinity
		foodMultiplier := float32(1.0) / (1.0 + overpopulationFactor*overpopulationFactor*2.0)
		food = rand.Float32() * foodMultiplier
	} else if populationRatio < 0.5 {
		// Low population (< 50%) - higher chance of max food
		// Use weighted random: 70% chance for high food (0.7-1.0), 30% for normal (0.0-1.0)
		if rand.Float32() < 0.7 {
			food = 0.7 + rand.Float32()*0.3 // 0.7 to 1.0
		} else {
			food = rand.Float32() // 0.0 to 1.0
		}
	} else {
		// Normal population (50% to 100%) - standard food
		food = rand.Float32()
	}

	c.Energy += food
	c.Energy -= 0.5

	// Death
	if c.Energy <= 0 {
		w.RemoveCreature(c.ID)
		return
	}

	// Reproduction
	if c.Energy >= w.ReproductionEnergyThreshold {
		// Try to find walkable position near parent
		for attempts := 0; attempts < 8; attempts++ {
			newPos := Vec2{
				X: c.Pos.X + rand.Intn(3) - 1,
				Y: c.Pos.Y + rand.Intn(3) - 1,
			}

			// Bounds check
			if newPos.X < 0 || newPos.Y < 0 || newPos.X >= w.Width || newPos.Y >= w.Height {
				continue
			}

			// Terrain check (1 = Plains = walkable)
			if w.TerrainMap[newPos.Y][newPos.X] == 1 {
				w.AddCreature(newPos, w.CreatureInitialEnergy)
				c.Energy -= w.ReproductionEnergyCost
				break
			}
		}
	}

	// Move
	newX := c.Pos.X + dx
	newY := c.Pos.Y + dy

	// Clamp to world bounds
	if newX < 0 {
		newX = 0
	}
	if newY < 0 {
		newY = 0
	}
	if newX >= w.Width {
		newX = w.Width - 1
	}
	if newY >= w.Height {
		newY = w.Height - 1
	}

	// Only move if terrain is walkable (1 = Plains)
	if w.TerrainMap[newY][newX] == 1 {
		c.Pos.X = newX
		c.Pos.Y = newY
	}
}
