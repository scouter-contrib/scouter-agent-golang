package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
)

//Ip4Value struct has number value
type Ip4Value struct {
	Value []byte
}

//NewIp4Value return DeciamlValue instance
func NewIp4Value(ip string) *Ip4Value {
	Ip4Value := new(Ip4Value)
	Ip4Value.Value = util.IpToBytes(ip)
	return Ip4Value
}

//NewIp4EmptyValue return Ip4Value instance
func NewIp4EmptyValue() *Ip4Value {
	Ip4Value := new(Ip4Value)
	Ip4Value.Value = make([]byte, 4)
	return Ip4Value
}

// Read function reads a value from datainputx
func (Ip4Value *Ip4Value) Read(in *DataInputX) (Value, error) {
	var err error
	Ip4Value.Value, err = in.Read(4)
	return Ip4Value, err
}

// Write function writes a number value to dataoutputx
func (Ip4Value *Ip4Value) Write(out *DataOutputX) error {
	err := out.Write(Ip4Value.Value)
	return err
}

// GetValueType returns value type
func (Ip4Value *Ip4Value) GetValueType() byte {
	return valueconstants.IP4ADDR
}

// ToString returns converted string value from decimal value
func (Ip4Value *Ip4Value) ToString() string {
	v := "ip["  + string(Ip4Value.Value) + "]"
	return v
}
