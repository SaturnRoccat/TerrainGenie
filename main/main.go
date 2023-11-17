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

	flag.IntVar(&gc.TotalChunksX, "X", 8, "The total number of chunks in the X axis")
	flag.IntVar(&gc.TotalChunksZ, "Z", 8, "The total number of chunks in the Z axis")
	flag.IntVar(&internals.ChunkHeight, "Y", 128, "The height of each chunk")

	flag.IntVar(&nc.Seed, "seed", rand.Int(), "The seed for the noise generator")
	flag.Float64Var(&nc.Frequency, "frequency", 0.01, "The frequency of the noise generator")
	flag.IntVar(&nc.Octaves, "octaves", 8, "The number of octaves for the noise generator")
	flag.Float64Var(&nc.Lacunarity, "lacunarity", 2.0, "The lacunarity of the noise generator")
	flag.Float64Var(&nc.Gain, "gain", 0.5, "The gain of the noise generator")
	flag.IntVar(&nc.Amplitude, "amplitude", 60, "The amplitude of the noise generator")

	flag.Parse()
	fmt.Println("Started terrain genie!...")

	if internals.ChunkHeight < nc.Amplitude {
		nc.Amplitude = internals.ChunkHeight - 3
	}

	ChunkArray := initChunkArray(&gc)
	wg := internals.MakeWG(nc)

	fillChunksWithHeightMap(&ChunkArray, &wg)

	fmt.Println("Done!")
}
