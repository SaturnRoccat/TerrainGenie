package main

import (
	"fmt"
	"os"
)

func addJsLine(TSCode *string, line string) {
	*TSCode += line + "\n"
}

func createJSDataNonRLE(config *TG_Config, palletData *TG_Pallet_Data) {
	var JSCode = "export const blockPallet = [" // Start of the pallet
	for _, blockName := range palletData.pallet {
		JSCode += "\"" + blockName + "\","
	}
	// remove trailing comma
	JSCode = JSCode[:len(JSCode)-1]
	JSCode += "];\n" // End of the pallet
	println(JSCode)
	JSCode += "import { Chunk } from \"./types.js\";\n"
	JSCode += "export const ChunkAmmountX = " + fmt.Sprint(config.XSize) + ";\n"
	JSCode += "export const ChunkAmmountZ = " + fmt.Sprint(config.ZSize) + ";\n"
	JSCode += "export const chunkSizeY = " + fmt.Sprint(ChunkHeight) + ";\n"
	JSCode += "export const chunkSizeXZ = " + fmt.Sprint(ChunkWidth) + ";\n"

	// Add the chunk data
	JSCode += "export const chunkData = [\n" // Start of the chunk data this could get thick
	for _, chunk := range WorldChunks {
		JSCode += "new Chunk(" + fmt.Sprint(chunk.ChunkPosition.x) + "," + fmt.Sprint(chunk.ChunkPosition.z)
		JSCode += "," + fmt.Sprint(ChunkWidth) + "," + fmt.Sprint(ChunkHeight)
		JSCode += ",["

		var byteArray [(ChunkWidth * ChunkDepth * ChunkHeight) * 4]byte
		overallIndex := 0
		for x := int32(0); x < ChunkWidth; x++ {
			for z := int32(0); z < ChunkDepth; z++ {
				for y := int32(0); y < ChunkHeight; y++ {
					var ourBytes = [4]byte{0x27, 0x61, 0x27, 0x2C}
					var index = TG_3D_PosToIndex(TG_3D_Pos{x, y, z})
					ourBytes[1] += byte(chunk.BlockData[index])

					byteArray[overallIndex*4] = ourBytes[0]
					byteArray[overallIndex*4+1] = ourBytes[1]
					byteArray[overallIndex*4+2] = ourBytes[2]
					byteArray[overallIndex*4+3] = ourBytes[3]
					overallIndex++
				}
			}
		}

		// Convert the byte array to a string
		conversionString := string(byteArray[:((ChunkWidth*ChunkDepth*ChunkHeight)*4)-1])
		JSCode += conversionString
		JSCode += "]),\n"
	}
	// remove trailing comma
	JSCode = JSCode[:len(JSCode)-2]
	JSCode += "];\n" // End of the chunk data

	os.WriteFile(config.JSOutputPath, []byte(JSCode), 0644)
}

func createJSDataRLE(config *TG_Config, palletData *TG_Pallet_Data) {
	var JSCode = "export const blockPallet = ["
	// remove trailing comma
	JSCode = JSCode[:len(JSCode)-1]
	JSCode += "];\n" // End of the pallet
	println(JSCode)

	// Add used data to the JSCode
	addJsLine(&JSCode, "export const chunkWidth = "+fmt.Sprint(config.XSize)+";")
	addJsLine(&JSCode, "export const chunkDepth = "+fmt.Sprint(config.ZSize)+";")
	addJsLine(&JSCode, "export const chunkHeight = "+fmt.Sprint(ChunkHeight)+";")
	addJsLine(&JSCode, "export const RLE = true;")
	print(JSCode)

	// Add the chunk data
	addJsLine(&JSCode, "export const chunkData = [\n") // Start of the chunk data this could get thick
	for _, chunk := range WorldChunks {
		addJsLine(&JSCode, "[")

		// Create RLE strings that represent the chunk
		for x := int32(0); x < ChunkWidth; x++ {
			for z := int32(0); z < ChunkDepth; z++ {
				var currentID uint16 = chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, 0, z})]
				var lowerY int32 = 0
				var UpperY int32 = 0
				for y := int32(0); y < int32(chunk.HeightMap[x+z*ChunkWidth]); y++ {
					if chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, y, z})] != currentID {
						// Add the RLE string
						addJsLine(&JSCode, "\""+fmt.Sprint(currentID)+"\":"+fmt.Sprint(lowerY)+":"+fmt.Sprint(UpperY)+",")
						// Update the RLE string
						currentID = chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{x, y, z})]
						println("CurrentID {} From: {} to {}", currentID, lowerY, UpperY)
						lowerY = y
						UpperY = y
					} else {
						UpperY++
					}
				}
			}
		}
	}

}
