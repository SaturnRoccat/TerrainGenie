package main

import (
	"image"
	"os"

	"golang.org/x/image/bmp"
)

func mapToFloat(max int, val int) float64 {
	if val < 0 || val > max {
		panic("Value is out of range")
	}

	return float64(val) / float64(max)
}

func bmpParseToChunk(cp TG_2D_Pos, imageData *image.Image, palletData *TG_Pallet_Data) {
	var chunk = makeTG_Level_Chunk(cp.x, cp.z)
	for x := 0; x < ChunkWidth; x++ {
		for z := 0; z < ChunkDepth; z++ {
			// Hacked in fix dont know if it works might just out of bounds
			var sx, sz = x * int(cp.x+1), z * int((cp.z + 1))
			var colorAtPixel = (*imageData).At(sx, sz)
			var r, _, _, _ = colorAtPixel.RGBA()
			//var height = uint16(mapToFloat(mv, int(r)) * float64(ChunkHeight))
			var height = uint16(mapToFloat(mv, int(r)) * float64(ChunkHeight))
			chunk.BlockData[TG_3D_PosToIndex(TG_3D_Pos{int32(x), int32(height), int32(z)})] = 1
			chunk.HeightMap[x+z*ChunkWidth] = height
		}
	}
	*(getFromWorldChunks(cp, chunkMapWidth)) = chunk
}

func importBMPHeightMap(config *TG_Config, palletData *TG_Pallet_Data) {
	// Load the BMP file
	var heightMap, err = os.Open(config.CustomHeightMap)
	if err != nil {
		panic(err)
	}
	// Decode the BMP file
	var imageData, anotherErr = bmp.Decode(heightMap)
	if anotherErr != nil {
		panic(anotherErr)
	}

	chunkMapWidth = int32(config.XSize)

	// calc if the r and g values are the same
	// Also find the max value
	var maxValue uint32 = 0
	var sameColor bool = true
	for x := 0; x < imageData.Bounds().Max.X; x++ {
		for z := 0; z < imageData.Bounds().Max.Y; z++ {
			var colorAtPixel = imageData.At(x, z)
			var r, g, _, _ = colorAtPixel.RGBA()
			if r != g {
				sameColor = false
				break
			}
			if r > maxValue {
				maxValue = r
			}
		}
	}
	mv = int(maxValue)
	if sameColor {
		println("The image is greyscale")
	} else {
		println("The image is not greyscale. Assuming that the red value is the height map value")
	}

	// Loop through the image and set the height map
	for x := 0; x < config.XSize; x++ {
		for z := 0; z < config.ZSize; z++ {
			bmpParseToChunk(TG_2D_Pos{int32(x), int32(z)}, &imageData, palletData)
		}
	}

}
