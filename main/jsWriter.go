package main

import (
	"duckos/TerrainGenie/main/internals"
	"os"
)

func exportTerrainDataToJS(oc *internals.OtherConfig, fillBuffer *internals.JSArrayInterface, chunkPositionBuffer *internals.JSArrayInterface, chunkSizeBuffer *internals.JSArrayInterface) {
	var jsData []byte
	jsData = append(jsData, (*fillBuffer).ToJSLine("commandBuff")...)
	jsData = append(jsData, (*chunkPositionBuffer).ToJSLine("chunkPositions")...)
	jsData = append(jsData, (*chunkSizeBuffer).ToJSLine("chunkSizes")...)

	var file, err = os.Create(oc.OutputJavasciptPath)
	if err != nil {
		panic(err)
	}
	file.Write(jsData)
	file.Close()
}
