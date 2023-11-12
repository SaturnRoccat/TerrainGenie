package main

import (
	"fmt"
)

var (
	WorldChunks = make([]TG_Level_Chunk, 0)
)

func getFromWorldChunks(pos TG_2D_Pos, XSize int32) *TG_Level_Chunk {
	return &WorldChunks[pos.x+pos.z*XSize]
}

func buildDataBuffer(config TG_Config, palletData *TG_Pallet_Data) {
	fmt.Println("Building data buffer...")
	WorldChunks = make([]TG_Level_Chunk, config.XSize*config.ZSize) // Allocate memory for the world we shouldnt need to do any more large allocations after this for the world generator

	var WorldGenerator = makeTG_Generator(config.Seed) // Make the world generator

	for x := int32(0); x < int32(config.XSize); x++ {
		for z := int32(0); z < int32(config.ZSize); z++ {
			var chunkData = getFromWorldChunks(TG_2D_Pos{x, z}, int32(config.XSize))
			var chunk = makeTG_Level_Chunk(x, z)
			blankedNoise(&chunk, &WorldGenerator)

			*chunkData = chunk
		}
	}
	fmt.Println("Building data buffer... Done!")

	turnWorldDataToBinary(&WorldChunks, &config, palletData)
}
