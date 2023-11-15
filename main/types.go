package main

// Path: main/types.go

import (
	"duckos/TerrainGenie/fastnoise"
	"encoding/binary"
	"math"
)

const (
	// Chunk dimensions
	ChunkWidth        = 16
	ChunkHeight       = 64
	ChunkDepth        = 16
	AmplitudeConstant = 50
)

type TG_Config struct {
	Seed                int
	OutputPath          string
	JSOutputPath        string
	XSize               int
	ZSize               int
	YSize               int
	EnableRLE           bool
	OutputNonCompressed bool
	CustomHeightMap     string
	// Noise parameters

	// Terrain blanket shape
	TerrainBlanketOctaves int
	TerrainBlanketLacun   float64
	TerrainBlanketGain    float64
	TerrainBlanketFreq    float64
	TerrainBlanketAmp     float64

	// CaveShape
	CaveShapeOctaves int
	CaveShapeLacun   float64
	CaveShapeGain    float64
	CaveShapeFreq    float64
	CaveShapeAmp     float64
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

type TG_Pallet_Data struct {
	pallet []string
}

type TG_2D_Pos TG_Chunk_Pos // Just another name for TG_Chunk_Pos with out having to redefine it

func TG_2D_PosToIndex(pos TG_2D_Pos) int32 {
	return pos.x + pos.z*ChunkWidth
}

func TG_3D_PosToIndex(pos TG_3D_Pos) int32 {
	return pos.x + pos.y*ChunkWidth + pos.z*ChunkWidth*ChunkHeight
}

func TG_Chunk_PosTo_3D_Pos(pos TG_Chunk_Pos) TG_3D_Pos {
	return TG_3D_Pos{pos.x*ChunkWidth - 1, 0, pos.z*ChunkDepth - 1}
}

func mapToZeroOne(value float32) float32 {
	// Ensure the value is within the valid range
	clampedValue := float32(math.Max(-1, math.Min(float64(value), 1)))

	// Map the value from the range [-1, 1] to [0, 1]
	mappedValue := (clampedValue + 1) / 2

	return mappedValue
}
func makeTG_Generator(config *TG_Config) TG_Generator {
	var gen TG_Generator
	gen.TerrainBlanketShapeState = fastnoise.New[float32]()
	gen.TerrainBlanketShapeState.Seed = config.Seed
	gen.TerrainBlanketShapeState.NoiseType(fastnoise.Perlin)
	gen.TerrainBlanketShapeState.Octaves = config.TerrainBlanketOctaves
	gen.TerrainBlanketShapeState.Lacunarity = float32(config.TerrainBlanketLacun)
	gen.TerrainBlanketShapeState.Gain = float32(config.TerrainBlanketGain)
	gen.TerrainBlanketShapeState.Frequency = float32(config.TerrainBlanketFreq)
	gen.TerrainBlanketShapeState.DomainWarpType = fastnoise.DomainWarpOpenSimplex2Reduced

	gen.CaveShapeState = fastnoise.New[float32]()
	gen.CaveShapeState.Seed = (((config.Seed * 897213) ^ config.Seed) * 219038) | config.Seed/2 // I don't know why this but its here now
	gen.CaveShapeState.NoiseType(fastnoise.Perlin)
	gen.CaveShapeState.Octaves = config.CaveShapeOctaves
	gen.CaveShapeState.Lacunarity = float32(config.CaveShapeLacun)
	gen.CaveShapeState.Gain = float32(config.CaveShapeGain)
	gen.CaveShapeState.Frequency = float32(config.CaveShapeFreq)

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
				if cave > 0.7 {
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

func addToPallet(pallet *TG_Pallet_Data, block string) uint16 {
	pallet.pallet = append(pallet.pallet, block)
	return uint16(len(pallet.pallet) - 1)
}

type uiint16 interface {
	uint16 | int16
}

type uiint8 interface {
	uint8 | int8
}

type uiint64 interface {
	uint64 | int64
}

type uiint32 interface {
	uint32 | int32
}

func int32ToBytes[T uiint32](value T) []byte {
	var buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(value))
	return buffer
}

func int16ToBytes[T uiint16](value T) []byte {
	var buffer = make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(value))
	return buffer
}

func int64ToBytes[T uiint64](value T) []byte {
	var buffer = make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(value))
	return buffer
}

func int8ToBytes[T uiint8](value T) []byte {
	var buffer = make([]byte, 1)
	buffer[0] = uint8(value)
	return buffer
}

func BytesToInt32[T uiint32](buffer []byte) T {
	return T(binary.BigEndian.Uint32(buffer))
}

func BytesToInt16[T uiint16](buffer []byte) T {
	return T(binary.BigEndian.Uint16(buffer))
}

func BytesToInt64[T uiint64](buffer []byte) T {
	return T(binary.BigEndian.Uint64(buffer))
}

func BytesToInt8[T uiint8](buffer []byte) T {
	return T(buffer[0])
}

func boolToUint8(value bool) uint8 {
	if value {
		return 1
	} else {
		return 0
	}
}
