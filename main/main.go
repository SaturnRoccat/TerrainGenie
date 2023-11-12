package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Welcome to Terrain Genie!")
	var config TG_Config

	flag.StringVar(&config.OutputPath, "output", "./levelData.bin", "Output path for generated terrain data")
	flag.IntVar(&config.Seed, "seed", rand.Int(), "Seed for random number generator")
	flag.Parse()

	fmt.Println("Output path:", config.OutputPath)
	fmt.Println("Seed:", config.Seed)

}
