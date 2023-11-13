package main

import (
	"fmt"
	"os"
)

func createSetblockCommand(arrayBuf *[]byte, x uint32, y uint32, z uint32, blockID uint16, palletData *TG_Pallet_Data) {
	*arrayBuf = append(*arrayBuf, []byte(fmt.Sprintf("\"setblock %d %d %d %s\",", x, y, z, palletData.pallet[blockID]))...)
}

func addJsLine(JSCode *string, line string) {
	*JSCode += line + "\n"
}

func createJSDataNonRLE(config *TG_Config, palletData *TG_Pallet_Data) {
	var outputJS []byte
	outputJS = append(outputJS, []byte("export const commandBuff = [")...)
	for _, chunk := range WorldChunks {
		for x := int32(0); x < ChunkWidth; x++ {
			for z := int32(0); z < ChunkDepth; z++ {
				for y := int32(0); y < ChunkHeight; y++ {
					var blockID = chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, y, z})]
					createSetblockCommand(&outputJS, uint32(x+chunk.ChunkPosition.x*ChunkWidth), uint32(y), uint32(z+chunk.ChunkPosition.z*ChunkDepth), blockID, palletData)
				}
			}
		}
	}
	outputJS[len(outputJS)-1] = ']' // Remove trailing comma
	outputJS = append(outputJS, []byte(";")...)
	os.WriteFile(config.JSOutputPath, outputJS, 0644)
}

func createJSDataRLE(config *TG_Config, palletData *TG_Pallet_Data) {

}
