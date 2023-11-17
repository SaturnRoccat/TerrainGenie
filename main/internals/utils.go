package internals

type OtherConfig struct {
	OutputJavascipt     bool
	OutputJavasciptPath string
}

func mapFloat(value float32) float32 {
	// Adjust the value to be between 0 and 2
	adjustedValue := (value + 1.0) / 2.0

	return adjustedValue
}
