package internals

import "fmt"

type JSArrayInterface interface {
	NewElement(element []byte)
	NewElementString(element string)
	ToString() string
	ToJSLine(exportName string) []byte
	Len() int
}

type JSArray struct {
	Length int
	Values [][]byte
	Raw    []byte
}

func NewJSArray() *JSArray {
	// Create a new JSArray
	var jsArray JSArray
	jsArray.Length = 0
	jsArray.Values = [][]byte{}
	jsArray.Raw = []byte{'['}
	return &jsArray
}

func (jsArray *JSArray) NewElement(element []byte) {
	jsArray.Values = append(jsArray.Values, element)
	jsArray.Raw = append(jsArray.Raw, element...)
	jsArray.Raw = append(jsArray.Raw, ',')
	jsArray.Length++
}

func (jsArray *JSArray) NewElementString(element string) {
	jsArray.Values = append(jsArray.Values, []byte(element))
	jsArray.Raw = append(jsArray.Raw, []byte(element)...)
	jsArray.Raw = append(jsArray.Raw, ',')
	jsArray.Length++
}

func (jsArray *JSArray) ToString() string {
	var str = string(jsArray.Raw)
	str = str[:len(str)-1] + "]"
	return str
}

func (jsArray *JSArray) ToJSLine(exportName string) []byte {
	return []byte(fmt.Sprintf("export const %s = %s;\n", exportName, jsArray.ToString()))
}

func (jsArray *JSArray) Len() int {
	return jsArray.Length
}
