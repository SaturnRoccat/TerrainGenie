package main

import (
	"fmt"
	"os"
	"unsafe"
)

func addJsLine(TSCode *string, line string) {
	*TSCode += line + "\n"
}

func createJSDataNonRLE(config *TG_Config, palletData *TG_Pallet_Data) {
	var TSCode = "export const blockPallet: string[] = [" // Start of the pallet
	for _, blockName := range palletData.pallet {
		TSCode += "\"" + blockName + "\","
	}
	// remove trailing comma
	TSCode = TSCode[:len(TSCode)-1]
	TSCode += "];\n" // End of the pallet
	print(TSCode)
	// Add used data to the TSCode
	addJsLine(&TSCode, "export const chunkWidth: number = "+fmt.Sprint(config.XSize)+";")
	addJsLine(&TSCode, "export const chunkDepth: number = "+fmt.Sprint(config.ZSize)+";")
	addJsLine(&TSCode, "export const chunkHeight: number = "+fmt.Sprint(ChunkHeight)+";")
	addJsLine(&TSCode, "export const RLE: boolean = false;")
	print(TSCode)
	// Add the chunk data
	addJsLine(&TSCode, "export const chunkData: string[][]  = [\n") // Start of the chunk data this could get thick
	for _, chunk := range WorldChunks {
		addJsLine(&TSCode, "[")

		var arrayAsBytes [(ChunkWidth * ChunkHeight * ChunkDepth) * 4]byte
		for index, block := range chunk.BlockData {
			ourBytes := [4]byte{0x22, 0x61, 0x22, 0x2C}
			ourBytes[1] += byte(block)
			arrayAsBytes[index*4] = ourBytes[0]
			arrayAsBytes[index*4+1] = ourBytes[1]
			arrayAsBytes[index*4+2] = ourBytes[2]
			arrayAsBytes[index*4+3] = ourBytes[3]
		}
		// Obtain a pointer to the array
		unsafeAABpointer := unsafe.Pointer(&arrayAsBytes)
		// convert that pointer to a byte pointer
		unsafeAABpointerAsBytePointer := (*byte)(unsafeAABpointer)

		addJsLine(&TSCode, unsafe.String(unsafeAABpointerAsBytePointer, ((ChunkWidth*ChunkHeight*ChunkDepth)*4)-1))
		addJsLine(&TSCode, "],")
		//print(TSCode)
	}
	TSCode = TSCode[:len(TSCode)-2]
	addJsLine(&TSCode, "];\n") // End of the chunk data

	// Write the file
	os.WriteFile(config.JSOutputPath, []byte(TSCode), 0644)
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
