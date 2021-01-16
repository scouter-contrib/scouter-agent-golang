package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//MapValue has map data
type MapValue struct {
	Table map[string]Value
}

//NewMapValue returns new MapValue instance
func NewMapValue() *MapValue {
	mapValue := new(MapValue)
	mapValue.Table = make(map[string]Value)
	return mapValue
}

// IsEmpty return whether mapvalue is empty or not
func (mapValue *MapValue) IsEmpty() bool {
	if len(mapValue.Table) == 0 {
		return true
	}
	return false
}

// ContainesKey returns whether mapvalue contains value for given key
func (mapValue *MapValue) ContainesKey(key string) bool {
	if _, ok := mapValue.Table[key]; ok {
		return true
	}
	return false

}

// Read will reads a value from datainputx
func (mapValue *MapValue) Read(in *DataInputX) (Value, error) {
	size, err := in.ReadDecimal()
	for i := int64(0); i < size; i++ {
		var key string
		var value Value
		key, err = in.ReadString()
		value, err = in.ReadValue()
		mapValue.Table[key] = value
	}
	return mapValue, err
}

// Write function write a text value to dataoutputx
func (mapValue *MapValue) Write(out *DataOutputX) error {
	_, err := out.WriteDecimal(int64(len(mapValue.Table)))
	for key, value := range mapValue.Table {
		_, err = out.WriteString(key)
		_, err = out.WriteValue(value)
	}
	return err
}

// GetValueType returns value type
func (mapValue *MapValue) GetValueType() byte {
	return valueconstants.MAP
}

// Put add string value to map
func (mapValue *MapValue) Put(key string, any interface{}) {
	switch any.(type) {
	case string:
		mapValue.Table[key] = NewTextValue(any.(string))
	case int64:
		mapValue.Table[key] = NewDecimalValue(any.(int64))
	case bool:
		mapValue.Table[key] = NewBooleanValue(any.(bool))
	case float32:
		mapValue.Table[key] = NewFloatValue(any.(float32))
	case Value:
		mapValue.Table[key] = any.(Value)
	}

}

// GetString returns string value
func (mapValue *MapValue) GetString(key string) string {
	stringValue, ok := mapValue.Table[key].(*TextValue)
	if ok {
		return stringValue.Value
	}
	return ""
}

// GetBoolean returns string value
func (mapValue *MapValue) GetBoolean(key string) bool {
	v, ok := mapValue.Table[key].(*BooleanValue)
	if ok {
		return v.Value
	}
	return false
}

// GetInt8 returns int8 value
func (mapValue *MapValue) GetInt8(key string) int8 {
	v, ok := mapValue.Table[key].(*DecimalValue)
	if ok {
		return int8(v.Value)
	}
	return 0
}

// GetInt16 returns int16 value
func (mapValue *MapValue) GetInt16(key string) int16 {
	v, ok := mapValue.Table[key].(*DecimalValue)
	if ok {
		return int16(v.Value)
	}
	return 0
}

// GetInt32 returns int32 value
func (mapValue *MapValue) GetInt32(key string) int32 {
	v, ok := mapValue.Table[key].(*DecimalValue)
	if ok {
		return int32(v.Value)
	}
	return 0
}

// GetInt64 returns int64 value
func (mapValue *MapValue) GetInt64(key string) int64 {
	v, ok := mapValue.Table[key].(*DecimalValue)
	if ok {
		return int64(v.Value)
	}
	return 0
}

// ToString returns converted string map data
func (mapValue *MapValue) ToString() string {
	str := "map value \n"
	for k, v := range mapValue.Table {
		str += ("key:" + k + " value:" + v.ToString() + "\n")
	}
	return str
}
