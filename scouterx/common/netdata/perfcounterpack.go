package netdata

import (
	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
)

const TimeTypeRealTime int8 = 1
const TimeTypeFiveMin int8 = 3

//PerfCounterPack has metric counter which has key/value type
type PerfCounterPack struct {
	Time     int64
	ObjName  string
	Timetype int8
	Data     *MapValue
}

// NewPerfCounterPack returns PerfCounterPack instance
func NewPerfCounterPack() *PerfCounterPack {
	pack := new(PerfCounterPack)
	pack.Data = NewMapValue()
	return pack
}

// Write will write PerfCounterPack to datoutputx
func (pack *PerfCounterPack) Write(out *DataOutputX) error {
	_, err := out.WriteInt64(pack.Time)
	_, err = out.WriteString(pack.ObjName)
	_, err = out.WriteInt8(pack.Timetype)
	_, err = out.WriteValue(pack.Data)
	return err
}

// Read will read PerfCounterPack from datainputx
func (pack *PerfCounterPack) Read(in *DataInputX) (Pack, error) {
	var err error
	pack.Time, err = in.ReadInt64()
	pack.ObjName, err = in.ReadString()
	pack.Timetype, err = in.ReadInt8()
	value, err := in.ReadValue()
	pack.Data = value.(*MapValue)
	return pack, err
}

// Put will put key/value to PerfCounterPack
func (pack *PerfCounterPack) Put(key string, any interface{}) {
	switch v := any.(type) {
	case int32:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case int64:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case uint32:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case uint64:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case int:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case float32:
		pack.Data.Put(key, NewFloatValue(v))
	case float64:
		pack.Data.Put(key, NewFloatValue(float32(v)))
	case string:
		pack.Data.Put(key, NewTextValue(v))
	case bool:
		pack.Data.Put(key, NewBooleanValue(v))
	default:
		pack.Data.Put(key, any)
	}
}

// ToString returns converted perfcounterpack value
func (pack *PerfCounterPack) ToString() string {
	str := pack.Data.ToString()
	return str
}

//GetPackType returns pack type
func (pack *PerfCounterPack) GetPackType() byte {
	return packconstants.PERFCOUNTER
}
