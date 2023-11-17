package main

import (
	"duckos/TerrainGenie/main/internals"
	"fmt"
)

func processCollumn(chunk *internals.Chunk, x, z int) []internals.RLECommand {
	var commands []internals.RLECommand

	// Start from the heightmap
	var height = chunk.HeightMap[x+z*internals.ChunkWidth]
	var currentBlockId = chunk.ChunkBlockData[internals.WorldPosition{X: x, Y: int(height), Z: z}.ToArrayPosition(internals.ChunkWidth, internals.ChunkHeight)]
	commands = append(commands, internals.RLECommand{Dir: internals.Down, StartPosition: internals.WorldPosition{X: x, Y: int(height), Z: z}, BlockID: currentBlockId, Length: 1}) // Add the first block
	for y := int(height); y > 0; y-- {
		var blockId = chunk.ChunkBlockData[internals.WorldPosition{X: x, Y: y, Z: z}.ToArrayPosition(internals.ChunkWidth, internals.ChunkHeight)]
		if blockId != currentBlockId {
			commands = append(commands, internals.RLECommand{Dir: internals.Down, StartPosition: internals.WorldPosition{X: x, Y: y, Z: z}, BlockID: currentBlockId, Length: 1})
			currentBlockId = blockId
		} else {
			commands[len(commands)-1].Length++
		}
	}

	return commands
}

func calculateRLEForChunk(chunk *internals.Chunk) []internals.RLECommand {
	var commands []internals.RLECommand
	for x := 0; x < internals.ChunkWidth; x++ {
		for z := 0; z < internals.ChunkWidth; z++ {
			commands = append(commands, processCollumn(chunk, x, z)...)
		}
	}
	return commands
}

func createFillBuffer(chunks *[]internals.Chunk, pallet *[]string) (internals.JSArrayInterface, internals.JSArrayInterface) {
	var fillBuffer internals.JSArrayInterface = internals.NewJSArray()
	var chunkSizeBuffer internals.JSArrayInterface = internals.NewJSArray()
	for _, chunk := range *chunks {
		var commands = calculateRLEForChunk(&chunk)
		for _, command := range commands {
			fillBuffer.NewElement(internals.RLECommandToCommand(command, pallet))
		}
		chunkSizeBuffer.NewElement([]byte(fmt.Sprintf("%d", len(commands))))
	}
	return fillBuffer, chunkSizeBuffer
}
