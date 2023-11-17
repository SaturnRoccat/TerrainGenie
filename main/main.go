package main

import (
	"duckos/TerrainGenie/main/internals"
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	var gc internals.GeneratorConfig
	var nc internals.NoiseConfig
	var oc internals.OtherConfig

	flag.IntVar(&gc.TotalChunksX, "X", 8, "The total number of chunks in the X axis")
	flag.IntVar(&gc.TotalChunksZ, "Z", 8, "The total number of chunks in the Z axis")
	flag.IntVar(&internals.ChunkHeight, "Y", 128, "The height of each chunk")

	flag.IntVar(&nc.Seed, "seed", rand.Int(), "The seed for the noise generator")
	flag.Float64Var(&nc.Frequency, "frequency", 0.01, "The frequency of the noise generator")
	flag.IntVar(&nc.Octaves, "octaves", 8, "The number of octaves for the noise generator")
	flag.Float64Var(&nc.Lacunarity, "lacunarity", 2.0, "The lacunarity of the noise generator")
	flag.Float64Var(&nc.Gain, "gain", 0.5, "The gain of the noise generator")
	flag.IntVar(&nc.Amplitude, "amplitude", 60, "The amplitude of the noise generator")

	flag.StringVar(&oc.OutputJavasciptPath, "output", "levelData.js", "The file to output the terrain data to")

	flag.Parse()
	fmt.Println("Started terrain genie!...")

	// We do this so that if in the rare case that the noise value generated is .9 before the conversion to int it will still be in the chunk and not out of bounds
	if internals.ChunkHeight < nc.Amplitude {
		nc.Amplitude = internals.ChunkHeight - 3
	}

	ChunkArray := initChunkArray(&gc)
	wg := internals.MakeWG(nc)

	fillChunksWithHeightMap(&ChunkArray, &wg)

	// Create the pallet
	var pallet []string
	pallet = append(pallet, "minecraft:air")
	pallet = append(pallet, "minecraft:stone")
	pallet = append(pallet, "minecraft:grass_block")

	// Create the fill buffer
	var fillBuffer, chunkSizeBuffer = createFillBuffer(&ChunkArray, &pallet)

	// Create the Chunk Positions buffer
	var chunkPositionBuffer = createChunkPositinBuffer(&ChunkArray)

	// Export the terrain data to JS
	exportTerrainDataToJS(&oc, &fillBuffer, &chunkPositionBuffer, &chunkSizeBuffer)

	fmt.Println("Done!")
}
