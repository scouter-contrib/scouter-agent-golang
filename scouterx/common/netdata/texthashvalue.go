package netdata

import (
	"strconv"

	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//TextHashValue struct has number value
type TextHashValue struct {
	Value int32
}

//NewTextHashValue return DeciamlValue instance
func NewTextHashValue(value int32) *TextHashValue {
	textHashValue := new(TextHashValue)
	textHashValue.Value = value
	return textHashValue
}

//NewDecimalEmptyValue return DeciamlValue instance
func NewTextHashEmptyValue() *TextHashValue {
	textHashValue := new(TextHashValue)
	return textHashValue
}

// Read function reads a value from datainputx
func (text *TextHashValue) Read(in *DataInputX) (Value, error) {
	var err error
	text.Value, err = in.ReadInt32()
	return text, err
}

// Write function writes a number value to dataoutputx
func (text *TextHashValue) Write(out *DataOutputX) error {
	_, err := out.WriteInt32(text.Value)
	return err
}

// GetValueType returns value type
func (text *TextHashValue) GetValueType() byte {
	return valueconstants.TEXT_HASH
}

// ToString returns converted string value from decimal value
func (text *TextHashValue) ToString() string {
	v := strconv.FormatInt(int64(text.Value), 10)
	return v
}
