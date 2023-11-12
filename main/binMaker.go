package main

import (
	"fmt"
	"os"
	"unsafe"
)

type TG_Header struct {
	WorldDataStart  int32
	PalletDataStart int32

	ChunkWidth  int16
	ChunkDepth  int16
	ChunkHeight int16
}

func int16arrayToBytes(array *[ChunkWidth * ChunkHeight * ChunkDepth]uint16) []byte {
	var buffer = make([]byte, 0)
	for _, value := range array {
		buffer = append(buffer, int16ToBytes(value)...)
	}
	return buffer
}

func addWorldDataToBinaryNoRLE(worldData *[]TG_Level_Chunk, dataBuffer *[]byte) {
	var chunkData = make([]byte, 0)
	for index, chunk := range *worldData {
		println("Adding Chunk:", index)
		chunkData = append(chunkData, int16arrayToBytes(&chunk.BlockData)...)
	}
	*dataBuffer = append(*dataBuffer, chunkData...)
}

func turnWorldDataToBinary(worldData *[]TG_Level_Chunk, config *TG_Config, palletData *TG_Pallet_Data) {
	var dataBuffer = make([]byte, 0)

	var header TG_Header

	calculateSizeOfPallet := func(pallet *TG_Pallet_Data) int32 {
		var palletSize = int32(0)
		for _, blockName := range pallet.pallet {
			palletSize += int32(unsafe.Sizeof(blockName))
		}
		return palletSize - int32(len(pallet.pallet)) // I shouldnt have to do this but for some reason if i dont do this the size is n bytes too big
	}

	addPalletToBuffer := func(pallet *TG_Pallet_Data) {
		for _, blockName := range pallet.pallet {
			dataBuffer = append(dataBuffer, []byte(blockName)...)
			dataBuffer = append(dataBuffer, 0x00)
		}
	}

	// Calculate the header
	header.PalletDataStart = int32(14) // 14 is the size of the header
	header.WorldDataStart = header.PalletDataStart + calculateSizeOfPallet(palletData)
	header.ChunkWidth = int16(config.XSize)
	header.ChunkHeight = int16(config.YSize)
	header.ChunkDepth = int16(config.ZSize)

	// Write the header
	dataBuffer = append(dataBuffer, int32ToBytes(header.WorldDataStart)...)
	dataBuffer = append(dataBuffer, int32ToBytes(header.PalletDataStart)...)
	dataBuffer = append(dataBuffer, int16ToBytes(header.ChunkWidth)...)
	dataBuffer = append(dataBuffer, int16ToBytes(header.ChunkDepth)...)
	dataBuffer = append(dataBuffer, int16ToBytes(header.ChunkHeight)...)
	// Print the header
	fmt.Print("Header:")
	fmt.Println(dataBuffer)

	// Write the pallet
	addPalletToBuffer(palletData)

	// write the RLE flag
	dataBuffer = append(dataBuffer, int8ToBytes(boolToUint8(config.EnableRLE))...)
	// Print the data Buffer
	fmt.Print("Data Buffer (This includes the header):")
	fmt.Println(dataBuffer)
	println("Added pallet to buffer")
	println("Adding world data to buffer...")
	if config.EnableRLE {
		// We cry here
	} else {
		addWorldDataToBinaryNoRLE(worldData, &dataBuffer)
	}
	println("Added world data to buffer saving to file...")
	permissions := 0644 // or whatever you need
	os.WriteFile(config.OutputPath, dataBuffer, os.FileMode(permissions))
	println("Saved to file! at path:", config.OutputPath)
}
