package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Welcome to Terrain Genie!")
	var config TG_Config
	var generateJS bool = false

	flag.StringVar(&config.OutputPath, "output", "./levelData.bin", "Output path for generated terrain data")
	flag.IntVar(&config.Seed, "seed", rand.Int(), "Seed for random number generator")
	flag.IntVar(&config.XSize, "X", 8, "The size in CHUNKS of the X axis of the world")
	flag.IntVar(&config.ZSize, "Z", 8, "The size in CHUNKS of the Z axis of the world")
	flag.IntVar(&config.YSize, "Y", 300, "The size in BLOCKS of the Y axis of the world")

	flag.Parse()

	fmt.Println("Output path:", config.OutputPath)
	fmt.Println("Seed:", config.Seed)

	if generateJS {
		fmt.Println("Generating JavaScript...")
	} else {
		fmt.Println("Generating binary of world data this will take some time...")
		buildDataBuffer(config)
	}

}
