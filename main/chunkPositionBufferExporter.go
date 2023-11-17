package main

import (
	"duckos/TerrainGenie/main/internals"
	"fmt"
)

func createChunkPositinBuffer(chunks *[]internals.Chunk) internals.JSArrayInterface {
	var chunkPositionBuffer internals.JSArrayInterface = internals.NewJSArray()
	for _, chunk := range *chunks {
		var CX = chunk.ChunkPosition.X
		var CZ = chunk.ChunkPosition.Z
		chunkPositionBuffer.NewElement([]byte(fmt.Sprintf("%d", CX)))
		chunkPositionBuffer.NewElement([]byte(fmt.Sprintf("%d", CZ)))
	}
	return chunkPositionBuffer
}
