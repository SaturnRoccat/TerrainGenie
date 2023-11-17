package internals

type Direction int

const (
	Up Direction = iota
	Down
)

type RLECommand struct {
	// The block ID
	BlockID uint16
	// The number of blocks
	Length int
	// The position of the first block
	StartPosition WorldPosition
	// The direction of the RLE
	Dir Direction
}

func handleSetBlock(rleData RLECommand, pallet *[]string) []byte {
	return makeSetBlockCommandString(rleData.StartPosition.X, rleData.StartPosition.Y, rleData.StartPosition.Z, rleData.BlockID, pallet)
}

func handleFill(rleData RLECommand, pallet *[]string) []byte {
	switch rleData.Dir {
	case Up:
		return makeFillCommandStringVer(rleData.StartPosition.X, rleData.StartPosition.Y, rleData.StartPosition.Z, rleData.StartPosition.X, rleData.StartPosition.Y+rleData.Length, rleData.StartPosition.Z, rleData.BlockID, pallet)
	case Down:
		return makeFillCommandStringVer(rleData.StartPosition.X, rleData.StartPosition.Y-rleData.Length, rleData.StartPosition.Z, rleData.StartPosition.X, rleData.StartPosition.Y, rleData.StartPosition.Z, rleData.BlockID, pallet)
	}
	return []byte{}
}

func RLECommandToCommand(rleData RLECommand, pallet *[]string) []byte {
	if rleData.Length == 1 {
		return handleSetBlock(rleData, pallet)
	} else {
		return handleFill(rleData, pallet)
	}
}
