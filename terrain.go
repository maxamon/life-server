package main

import (
	"math"
	"math/rand"
)

type TerrainType byte

const (
	Water TerrainType = iota
	Plains
	Mountains
)

type TerrainGenerator struct {
	Seed         int64
	WorldWidth   int
	WorldHeight  int
	BorderSize   int
	permutation  [512]int
	gradients    [][2]float64
}

func NewTerrainGenerator(seed int64, width, height int) *TerrainGenerator {
	tg := &TerrainGenerator{
		Seed:        seed,
		WorldWidth:  width,
		WorldHeight: height,
		BorderSize:  50, // Water border width
	}

	// Initialize permutation table with seed
	rng := rand.New(rand.NewSource(seed))
	perm := rng.Perm(256)
	for i := 0; i < 256; i++ {
		tg.permutation[i] = perm[i]
		tg.permutation[256+i] = perm[i]
	}

	// Generate gradient vectors
	tg.gradients = make([][2]float64, 256)
	for i := 0; i < 256; i++ {
		angle := rng.Float64() * 2 * math.Pi
		tg.gradients[i] = [2]float64{math.Cos(angle), math.Sin(angle)}
	}

	return tg
}

func (tg *TerrainGenerator) GetTerrain(x, y int) TerrainType {
	// Force water at borders
	if x < tg.BorderSize || y < tg.BorderSize ||
		x >= tg.WorldWidth-tg.BorderSize || y >= tg.WorldHeight-tg.BorderSize {
		return Water
	}

	// Perlin noise
	noise := tg.perlin(float64(x)*0.01, float64(y)*0.01)

	// Map noise [-1, 1] to terrain types
	if noise < -0.2 {
		return Water
	} else if noise < 0.5 {
		return Plains
	}
	return Mountains
}

func (tg *TerrainGenerator) perlin(x, y float64) float64 {
	// Grid cell coordinates
	xi := int(math.Floor(x)) & 255
	yi := int(math.Floor(y)) & 255

	// Fractional part
	xf := x - math.Floor(x)
	yf := y - math.Floor(y)

	// Fade curves
	u := fade(xf)
	v := fade(yf)

	// Hash coordinates of 4 corners
	aa := tg.permutation[tg.permutation[xi]+yi]
	ab := tg.permutation[tg.permutation[xi]+yi+1]
	ba := tg.permutation[tg.permutation[xi+1]+yi]
	bb := tg.permutation[tg.permutation[xi+1]+yi+1]

	// Gradients at corners
	g1 := dot(tg.gradients[aa], xf, yf)
	g2 := dot(tg.gradients[ba], xf-1, yf)
	g3 := dot(tg.gradients[ab], xf, yf-1)
	g4 := dot(tg.gradients[bb], xf-1, yf-1)

	// Interpolate
	x1 := lerp(u, g1, g2)
	x2 := lerp(u, g3, g4)

	return lerp(v, x1, x2)
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func dot(g [2]float64, x, y float64) float64 {
	return g[0]*x + g[1]*y
}

func (tg *TerrainGenerator) IsWalkable(x, y int) bool {
	return tg.GetTerrain(x, y) == Plains
}
