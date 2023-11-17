package internals

import (
	"duckos/TerrainGenie/fastnoise"
)

type GeneratorConfig struct {
	TotalChunksX, TotalChunksZ int
}

type NoiseConfig struct {
	Seed       int
	Frequency  float64
	Octaves    int
	Lacunarity float64
	Gain       float64
	Nt         fastnoise.NoiseType
	Amplitude  int
}

type WorldGeneratorInterface interface {
	GetHeightMap(x, z int, writeMap *[ChunkWidth * ChunkWidth]uint16)
	GetHeightAt(x, z int) int
}

type WorldGenerator struct {
	TerrainShape2D *fastnoise.State[float32]
	amplitude      int
	config         NoiseConfig
}

// Generates a heightmap for a chunk
func (wg WorldGenerator) GetHeightMap(chunkX, chunkZ int, writeMap *[ChunkWidth * ChunkWidth]uint16) {
	for lx := 0; lx < ChunkWidth; lx++ {
		x := lx + (chunkX * ChunkWidth)
		for lz := 0; lz < ChunkWidth; lz++ {
			z := lz + (chunkZ * ChunkWidth)
			// Ensure GetHeightAt returns non-zero values for x and z coordinates
			height := uint16(wg.GetHeightAt(x, z))
			// Calculate the index properly
			index := lx + lz*ChunkWidth
			// Assign the height value to the writeMap array
			(*writeMap)[index] = height
		}
	}
}

func (wg WorldGenerator) GetHeightAt(x, z int) int {
	return int(mapFloat(wg.TerrainShape2D.GetNoise2D(float32(x), float32(z))) * float32(wg.amplitude))
}

func MakeWG(config NoiseConfig) WorldGenerator {
	var wg WorldGenerator

	wg.config = config
	wg.amplitude = config.Amplitude
	// Init the terrain shape state
	wg.TerrainShape2D = fastnoise.New[float32]()
	wg.TerrainShape2D.NoiseType(fastnoise.OpenSimplex2S)
	wg.TerrainShape2D.Seed = config.Seed
	wg.TerrainShape2D.Octaves = config.Octaves
	wg.TerrainShape2D.Frequency = float32(config.Frequency)
	wg.TerrainShape2D.Lacunarity = float32(config.Lacunarity)
	wg.TerrainShape2D.Gain = float32(config.Gain)
	return wg
}
