package main

import (
	"bytes"
	"compress/gzip"
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

func compressedint16arrayToBytes(array *[ChunkWidth * ChunkHeight * ChunkDepth]uint16) []byte {
	var buffer = int16arrayToBytes(array)
	var compressedBuffer bytes.Buffer
	var writer = gzip.NewWriter(&compressedBuffer)
	_, err := writer.Write(buffer)
	if err != nil {
		return nil
	}

	return compressedBuffer.Bytes()
}

func addWorldDataToBinaryNoRLE(worldData *[]TG_Level_Chunk, dataBuffer *[]byte) {
	var chunkData []byte
	for index, chunk := range *worldData {
		println("Adding Chunk:", index)
		chunkData = compressedint16arrayToBytes(&chunk.BlockData)
		*dataBuffer = append(*dataBuffer, 0xFF, 0xAA, 0xFF, 0xAA) // Chunk start flag
		*dataBuffer = append(*dataBuffer, chunkData...)
		*dataBuffer = append(*dataBuffer, 0xAA, 0xFF, 0xAA, 0xFF) // Chunk end flag
	}
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

	var buffer bytes.Buffer
	var writer = gzip.NewWriter(&buffer)

	_, err := writer.Write(dataBuffer)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(config.OutputPath+".cbin", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	buffer.Reset()

	if config.OutputNonCompressed {
		os.WriteFile(config.OutputPath+".ubin", dataBuffer, 0644)
	}

	println("Saved to file! at path:", config.OutputPath)
}
