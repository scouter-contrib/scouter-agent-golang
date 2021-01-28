package netdata

import (
	"strconv"

	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

type DoubleValue struct {
	Value float64
}

func NewDoubleValue(value float64) *DoubleValue {
	doubleValue := new(DoubleValue)
	doubleValue.Value = value
	return doubleValue
}

//NewFloatEmptyValue returns new FloatValue instance
func NewDoubleEmptyValue() *DoubleValue {
	doubleValue := new(DoubleValue)
	return doubleValue
}

// Read function reads a value from datainputx
func (value *DoubleValue) Read(in *DataInputX) (Value, error) {
	var err error
	value.Value, err = in.ReadFloat64()
	return value, err
}

// Write function write a float value to dataoutputx
func (value *DoubleValue) Write(out *DataOutputX) error {
	_, err := out.WriteFloat64(value.Value)
	return err
}

// GetValueType returns value type
func (value *DoubleValue) GetValueType() byte {
	return valueconstants.DOUBLE
}

// ToString returns converted float value
func (value *DoubleValue) ToString() string {
	return strconv.FormatFloat(float64(value.Value), 'f', 3, 64)
}
