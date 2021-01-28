package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
	"strconv"
)

//BlobValue struct has number value
type BlobValue struct {
	Value []byte
}

//NewBlobValue return DeciamlValue instance
func NewBlobValue(value []byte) *BlobValue {
	BlobValue := new(BlobValue)
	BlobValue.Value = value
	return BlobValue
}

//NewBlobEmptyValue return BlobValue instance
func NewBlobEmptyValue() *BlobValue {
	BlobValue := new(BlobValue)
	return BlobValue
}

// Read function reads a value from datainputx
func (BlobValue *BlobValue) Read(in *DataInputX) (Value, error) {
	var err error
	BlobValue.Value, err = in.ReadBlob()
	return BlobValue, err
}

// Write function writes a number value to dataoutputx
func (BlobValue *BlobValue) Write(out *DataOutputX) error {
	err := out.WriteBlob(BlobValue.Value)
	return err
}

// GetValueType returns value type
func (BlobValue *BlobValue) GetValueType() byte {
	return valueconstants.BLOB
}

// ToString returns converted string value from decimal value
func (BlobValue *BlobValue) ToString() string {
	v := "byte["  + strconv.Itoa(len(BlobValue.Value)) + "]"
	return v
}
