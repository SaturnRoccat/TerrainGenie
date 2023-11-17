package main

import (
	"duckos/TerrainGenie/main/internals"
)

// Fills the chunk with the heightmap of that chunk! This is a very simple implementation and will be replaced with a more complex one later
func fillChunksWithHeightMap(chunks *[]internals.Chunk, wgi internals.WorldGeneratorInterface) {
	for _, chunk := range *chunks {
		wgi.GetHeightMap(int(chunk.ChunkPosition.X), int(chunk.ChunkPosition.Z), &chunk.HeightMap)

		// Merge the heightmap into the chunk
		for x := 0; x < internals.ChunkWidth; x++ {
			for z := 0; z < internals.ChunkWidth; z++ {
				var height = chunk.HeightMap[x+z*internals.ChunkWidth]
				var index = internals.WorldPosition{X: x, Y: int(height), Z: z}.ToArrayPosition(internals.ChunkWidth, internals.ChunkHeight)

				chunk.ChunkBlockData[index] = 2 // The id of grass. This needs to get reworked for biome support
				for y := 0; y < int(height-2); y++ {
					chunk.ChunkBlockData[internals.WorldPosition{X: x, Y: y, Z: z}.ToArrayPosition(internals.ChunkWidth, internals.ChunkHeight)] = 1 // The id of stone
				}
			}
		}
	}
}

// Inits the chunks
func initChunkArray(config *internals.GeneratorConfig) []internals.Chunk {
	var chunks []internals.Chunk
	for x := 0; x < config.TotalChunksX; x++ {
		for z := 0; z < config.TotalChunksZ; z++ {
			chunks = append(chunks, internals.InitChunk(x, z, internals.ChunkHeight))
		}
	}
	return chunks
}
