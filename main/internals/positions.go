package internals

// Path: main/internals/positions.go

type WorldPositionInterface interface {
	ToArrayPosition() int
}

type WorldPosition struct {
	X, Y, Z int
}

func (pos WorldPosition) ToArrayPosition(chunkWidth, chunkHeight int) int {
	return pos.X + pos.Y*chunkWidth + pos.Z*chunkWidth*chunkHeight
}

type TwoDPosition struct {
	X, Z int
}

func (pos TwoDPosition) ToArrayPosition(chunkWidth int) int {
	return pos.X + pos.Z*chunkWidth
}
