package main

// Path: main/types.go

import (
	"duckos/TerrainGenie/fastnoise"
	"math"
)

const (
	// Chunk dimensions
	ChunkWidth        = 16
	ChunkHeight       = 358
	ChunkDepth        = 16
	AmplitudeConstant = 120
)

type TG_Config struct {
	Seed       int
	OutputPath string
	XSize      int
	ZSize      int
	YSize      int
	// Noise parameters
}

type TG_Generator struct {
	// Basic 2D shape noise
	TerrainBlanketShapeState *fastnoise.State[float32]

	// We use float32 because it's smaller than float64 and we don't need the precision also we need alot of these

	// basic 3D noise for caves
	CaveShapeState *fastnoise.State[float32]
}

type TG_Chunk_Pos struct {
	x, z int32
}

type TG_3D_Pos struct {
	x, y, z int32
}

type TG_Level_Chunk struct {
	BlockData     [ChunkWidth * ChunkHeight * ChunkDepth]uint16
	ChunkPosition TG_Chunk_Pos // Chunk position in chunk coordinates

	HeightMap [ChunkWidth * ChunkWidth]uint16
}

type TG_2D_Pos TG_Chunk_Pos // Just another name for TG_Chunk_Pos with out having to redefine it

func TG_3D_PosToIndex(pos TG_3D_Pos) int32 {
	return pos.x + pos.y*ChunkWidth + pos.z*ChunkWidth*ChunkWidth
}

func TG_Chunk_PosTo_3D_Pos(pos TG_Chunk_Pos) TG_3D_Pos {
	return TG_3D_Pos{pos.x * ChunkWidth, 0, pos.z * ChunkDepth}
}

func mapToZeroOne(value float32) float32 {
	// Ensure the value is within the valid range
	clampedValue := float32(math.Max(-1, math.Min(float64(value), 1)))

	// Map the value from the range [-1, 1] to [0, 1]
	mappedValue := (clampedValue + 1) / 2

	return mappedValue
}
func makeTG_Generator(seed int) TG_Generator {
	var gen TG_Generator
	gen.TerrainBlanketShapeState = fastnoise.New[float32]()
	gen.TerrainBlanketShapeState.Seed = seed
	gen.TerrainBlanketShapeState.NoiseType(fastnoise.OpenSimplex2S)

	gen.CaveShapeState = fastnoise.New[float32]()
	gen.CaveShapeState.Seed = (((seed * 897213) ^ seed) * 219038) | seed/2 // I don't know why this but its here now
	gen.CaveShapeState.NoiseType(fastnoise.OpenSimplex2S)
	return gen
}

func blankedNoise(chunk_p *TG_Level_Chunk, generator_p *TG_Generator) {
	var ChunkPAsWorldP = TG_Chunk_PosTo_3D_Pos(chunk_p.ChunkPosition)

	// Creates a 2D heightmap for the chunk
	for x := int32(0); x < ChunkWidth; x++ {
		var WorldP = x + ChunkPAsWorldP.x
		for z := int32(0); z < ChunkDepth; z++ {
			var WorldZ = z + ChunkPAsWorldP.z
			var height = int32(
				mapToZeroOne(generator_p.TerrainBlanketShapeState.GetNoise2D(float32(WorldP), float32(WorldZ))) * AmplitudeConstant) // Get height from noise this should be illigal and it should be in a lambda or function but I'm lazy
			chunk_p.HeightMap[x+z*ChunkWidth] = uint16(height)
			chunk_p.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, height, z})] = 1 // Temp constant for stone this should be replaced by a cached block id
		}
	}

	// Set everything below the heightmap to stone
	for x := int32(0); x < ChunkWidth; x++ {
		for z := int32(0); z < ChunkDepth; z++ {
			for y := int32(0); y < int32(chunk_p.HeightMap[x+z*ChunkWidth]); y++ {
				chunk_p.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, y, z})] = 1
			}
		}
	}

	// Create caves
	for x := int32(0); x < ChunkWidth; x++ {
		var WorldP = x + ChunkPAsWorldP.x
		for z := int32(0); z < ChunkDepth; z++ {
			var WorldZ = z + ChunkPAsWorldP.z
			for y := int32(0); y < int32(chunk_p.HeightMap[x+z*ChunkWidth]); y++ {
				var cave = mapToZeroOne(generator_p.CaveShapeState.GetNoise3D(float32(WorldP), float32(y), float32(WorldZ)))
				if cave > 0.5 {
					chunk_p.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, y, z})] = 0
				}
			}
		}
	}
}

func makeTG_Level_Chunk(x, z int32) TG_Level_Chunk {
	var chunk TG_Level_Chunk
	chunk.ChunkPosition.x = x
	chunk.ChunkPosition.z = z
	return chunk
}
