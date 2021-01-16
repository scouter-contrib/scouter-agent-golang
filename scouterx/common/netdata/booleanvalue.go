package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//BooleanValue has bool value
type BooleanValue struct {
	Value bool
}

//NewBooleanValue return new BooleanVaue instance
/*
func NewBooleanValue(value bool) *BooleanValue {
	booleanValue := new(BooleanValue)
	booleanValue.Value = value
	return booleanValue

}
*/

//NewBooleanValue return new BooleanVaue instance
func NewBooleanValue(value bool) *BooleanValue {
	booleanValue := new(BooleanValue)
	booleanValue.Value = value
	return booleanValue

}

//NewBooleanEmptyValue return new BooleanVaue instance
func NewBooleanEmptyValue() *BooleanValue {
	BooleanValue := new(BooleanValue)
	return BooleanValue
}

//GetValueType returns value type
func (booleanValue *BooleanValue) GetValueType() byte {
	return valueconstants.BOOLEAN
}

func (booleanValue *BooleanValue) Read(in *DataInputX) (Value, error) {
	var err error
	booleanValue.Value, err = in.ReadBoolean()
	return booleanValue, err
}

func (booleanValue *BooleanValue) Write(out *DataOutputX) error {
	_, err := out.WriteBoolean(booleanValue.Value)
	return err
}

// ToString returns converted string value from boolean value
func (booleanValue *BooleanValue) ToString() string {
	if booleanValue.Value {
		return "true"
	}
	return "fasle"
}
