package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Welcome to Terrain Genie!")
	var config TG_Config
	var generateJS bool = true

	flag.StringVar(&config.OutputPath, "output", "./levelData", "Output path for generated terrain data")
	flag.StringVar(&config.JSOutputPath, "TSOutput", "./levelData.ts", "Output path for the generated JavaScript file (This is only used if the -TS flag is set)")
	flag.IntVar(&config.Seed, "seed", rand.Int(), "Seed for random number generator")
	flag.IntVar(&config.XSize, "X", 24, "The size in CHUNKS of the X axis of the world")
	flag.IntVar(&config.ZSize, "Z", 24, "The size in CHUNKS of the Z axis of the world")
	flag.IntVar(&config.YSize, "Y", 300, "The size in BLOCKS of the Y axis of the world")
	flag.BoolVar(&config.OutputNonCompressed, "ONC", false, "Ouput non compressed binary data")
	flag.BoolVar(&config.EnableRLE, "RLE", false, "Enable RLE compression this reduces memory usage but takes longer to generate")
	flag.BoolVar(&generateJS, "TS", true, "Generate Typesscript from the world data")

	flag.Parse()

	fmt.Println("Output path:", config.OutputPath)
	fmt.Println("Seed:", config.Seed)

	// Pallet Data
	var palletData = TG_Pallet_Data{}
	addToPallet(&palletData, "minecraft:air")
	addToPallet(&palletData, "minecraft:stone")

	fmt.Println("Generating binary of world data this will take some time...")
	buildDataBuffer(config, &palletData)

	if generateJS {
		fmt.Println("Generating JavaScript...")
		createJSData(&config, &palletData)
	}
}
