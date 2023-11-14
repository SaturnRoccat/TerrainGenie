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
	var generateHeightMapVer bool = false

	flag.StringVar(&config.OutputPath, "output", "./levelData", "Output path for generated terrain data")
	flag.StringVar(&config.JSOutputPath, "TSOutput", "./levelData.js", "Output path for the generated JavaScript file (This is only used if the -JS flag is set)")
	flag.IntVar(&config.Seed, "seed", rand.Int(), "Seed for random number generator")
	flag.IntVar(&config.XSize, "X", 4, "The size in CHUNKS of the X axis of the world")
	flag.IntVar(&config.ZSize, "Z", 4, "The size in CHUNKS of the Z axis of the world")
	flag.IntVar(&config.YSize, "Y", 256, "The size in BLOCKS of the Y axis of the world")
	flag.BoolVar(&config.OutputNonCompressed, "ONC", true, "Output non compressed binary data")
	flag.BoolVar(&config.EnableRLE, "RLE", false, "This is experimental and may not work")
	flag.BoolVar(&generateJS, "JS", true, "Generate JavaScript from the world data")
	flag.BoolVar(&generateHeightMapVer, "GMV", true, "Generate JavaScript from the world data only exporting the height map this makes smaller filles but not caves")
	config.Seed = 10

	flag.IntVar(&config.TerrainBlanketOctaves, "TBO", 6, "Terrain blanket octaves")
	flag.Float64Var(&config.TerrainBlanketLacun, "TBL", 1.8, "Terrain blanket lacunarity")
	flag.Float64Var(&config.TerrainBlanketGain, "TBG", 0.6, "Terrain blanket gain")
	flag.Float64Var(&config.TerrainBlanketFreq, "TBF", 0.02, "Terrain blanket frequency")

	flag.IntVar(&config.CaveShapeOctaves, "CSO", 4, "Cave shape octaves")
	flag.Float64Var(&config.CaveShapeLacun, "CSL", 2.0, "Cave shape lacunarity")
	flag.Float64Var(&config.CaveShapeGain, "CSG", 0.5, "Cave shape gain")
	flag.Float64Var(&config.CaveShapeFreq, "CSF", 0.01, "Cave shape frequency")
	flag.Parse()

	fmt.Println("Output path:", config.OutputPath)
	fmt.Println("Seed:", config.Seed)

	// Pallet Data
	var palletData = TG_Pallet_Data{}
	addToPallet(&palletData, "minecraft:air")
	addToPallet(&palletData, "minecraft:stone")

	fmt.Println("Generating binary of world data this will take some time...")
	buildDataBuffer(config, &palletData)

	// THIS MUST ALLWAYS BE THE LAST THING TO RUN
	// The reason for this is because it uses unsafe pointers and it modifies the data in the world chunks
	if generateJS {
		fmt.Println("Generating JavaScript...")
		if config.EnableRLE {
			createJSDataRLE(&config, &palletData)
		} else {
			if generateHeightMapVer {
				createJSHeightMapVersion(&config, &palletData)
			} else {
				createJSDataNonRLE(&config, &palletData)
			}
		}
	}
}
