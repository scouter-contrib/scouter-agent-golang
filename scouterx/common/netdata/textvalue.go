package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

// TextValue contains string value
type TextValue struct {
	Value string
}

// NewTextValue returns TextValue instance
func NewTextValue(value string) *TextValue {
	textValue := new(TextValue)
	textValue.Value = value
	return textValue
}

// NewTextEmptyValue returns TextValue instance
func NewTextEmptyValue() *TextValue {
	textValue := new(TextValue)
	return textValue
}

// Read function reads a value from datainputx
func (textValue *TextValue) Read(in *DataInputX) (Value, error) {
	var err error
	textValue.Value, err = in.ReadString()
	return textValue, err
}

// Write function write a text value to dataoutputx
func (textValue *TextValue) Write(out *DataOutputX) error {
	_, err := out.WriteString(textValue.Value)
	return err
}

// GetValueType returns value type
func (textValue *TextValue) GetValueType() byte {
	return valueconstants.TEXT
}

// ToString returns value
func (textValue *TextValue) ToString() string {
	return textValue.Value
}
