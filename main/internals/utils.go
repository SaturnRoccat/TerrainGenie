package internals

import (
	"fmt"
)

type OtherConfig struct {
	OutputJavascipt     bool
	OutputJavasciptPath string
}

func mapFloat(value float32) float32 {
	// Adjust the value to be between 0 and 2
	adjustedValue := (value + 1.0) / 2.0

	return adjustedValue
}

/*
The reason we are using byte arrays instead of strings is for the speed increase.
Also strings are immutable, so we would have to create a new string every time we wanted to append to it.
Which would be very slow. So doing it this way is much faster. :D (I think)
*/

func makeSetBlockCommandString(x, y, z int, blockID uint16, pallet *[]string) []byte {
	return []byte(fmt.Sprintf("\"setblock %d %d %d %s\",", x, y, z, (*pallet)[blockID]))
}

func makeFillCommandStringVer(x1, y1, z1, x2, y2, z2 int, blockID uint16, pallet *[]string) []byte {
	return []byte(fmt.Sprintf("\"fill %d %d %d %d %d %d %s\",", x1, y1, z1, x2, y2, z2, (*pallet)[blockID]))
}
