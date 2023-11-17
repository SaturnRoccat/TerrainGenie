package internals

// "duckos/TerrainGenie/fastnoise"

const (
	ChunkWidth = 16
)

var (
	ChunkHeight = 128
)

func InitChunk(x, z, chunkHeight int) Chunk {
	var chunk Chunk
	chunk.ChunkPosition = TwoDPosition{x, z}
	chunk.ChunkBlockData = make([]uint16, ChunkWidth*chunkHeight*ChunkWidth)
	return chunk
}

type Chunk struct {
	// The chunk's position in the world
	ChunkPosition TwoDPosition
	// The blocks in the chunk
	ChunkBlockData []uint16

	// The Heightmap
	HeightMap [ChunkWidth * ChunkWidth]uint16
}

func BlockIndex(wp WorldPosition) int {
	return wp.ToArrayPosition(ChunkWidth, ChunkHeight)
}
