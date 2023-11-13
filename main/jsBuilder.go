package main

import (
	"fmt"
	"reflect"
	"unsafe" // We seem to be using this alot lmao
)

func addJsLine(TSCode *string, line string) {
	*TSCode += line + "\n"
}

func createJSData(config *TG_Config, palletData *TG_Pallet_Data) {
	var TSCode = "export const: string[] blockPallet = [" // Start of the pallet
	for _, blockName := range palletData.pallet {
		TSCode += "\"" + blockName + "\","
	}
	// remove trailing comma
	TSCode = TSCode[:len(TSCode)-1]
	TSCode += "];\n" // End of the pallet
	print(TSCode)
	// Add used data to the TSCode
	addJsLine(&TSCode, "export const: number chunkWidth = "+fmt.Sprint(config.XSize)+";")
	addJsLine(&TSCode, "export const: number chunkDepth = "+fmt.Sprint(config.ZSize)+";")
	addJsLine(&TSCode, "export const: number chunkHeight = "+fmt.Sprint(ChunkHeight)+";")
	print(TSCode)
	// Add the chunk data
	addJsLine(&TSCode, "export const: string[][] chunkData = [\n") // Start of the chunk data this could get thick
	for _, chunk := range WorldChunks {
		addJsLine(&TSCode, "[")

		// Point to the start of the chunk data in memory
		// This is risky but it should work
		// It is also very fast
		// We can just tell TS to read the string as hex
		// This is a very hacky way of doing this but it should work
		PointerToChunkData := uintptr(unsafe.Pointer(&chunk.BlockData[0]))
		// Convert the pointer to a byte array
		ResultingByteArray := (*[(ChunkWidth * ChunkHeight * ChunkDepth) * 2]byte)(unsafe.Pointer(PointerToChunkData))

		// Set every second byte to 0x2C (,) this is the seperator for the hex values
		for i := 1; i < len(ResultingByteArray); i += 2 {
			ResultingByteArray[i] = 0x2C
		}
		// Add 33 to every other byte this is to make the hex values valid and can be shown as a string
		for i := 0; i < len(ResultingByteArray); i += 2 {
			ResultingByteArray[i] += 33
		}

		PointerToByteArray := uintptr(unsafe.Pointer(&ResultingByteArray[0]))
		stringHeader := &reflect.StringHeader{
			Data: PointerToByteArray,
			Len:  (ChunkWidth * ChunkHeight * ChunkDepth) * 2,
		}
		resultingString := *(*string)(unsafe.Pointer(stringHeader))
		addJsLine(&TSCode, resultingString)
		addJsLine(&TSCode, "],")
		// print(TSCode)
	}
	addJsLine(&TSCode, "];\n") // End of the chunk data
	print(TSCode)
}
