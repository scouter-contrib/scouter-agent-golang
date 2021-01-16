package netdata

import (
	"strconv"

	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//ListValue has Value type Slice
type ListValue struct {
	Value []Value
}

//NewListValue returns new ListVaue instance
func NewListValue() *ListValue {
	listValue := new(ListValue)
	listValue.Value = make([]Value, 0)
	return listValue
}

/*
//NewListValueWithSize returns new ListVaue instance
func NewListValueWithSize(size int) *ListValue {
	listValue := new(ListValue)
	listValue.Value = make([]Value, size, size)
	return listValue
}
*/

//NewListValueWithValue returns new ListVaue instance
func NewListValueWithValue(value []Value) *ListValue {
	listValue := new(ListValue)
	listValue.Value = value
	return listValue
}

// Add function adds a value to slice
func (listValue *ListValue) Add(value Value) *ListValue {
	listValue.Value = append(listValue.Value, value)
	return listValue
}

// AddInt64 adds a int64 value to slice
func (listValue *ListValue) AddInt64(value int64) *ListValue {
	listValue.Value = append(listValue.Value, NewDecimalValue(value))
	return listValue
}

// AddInt32 adds a int64 value to slice
func (listValue *ListValue) AddInt32(value int32) *ListValue {
	listValue.Value = append(listValue.Value, NewDecimalValue(int64(value)))
	return listValue
}

// AddFloat add a float value to slice
func (listValue *ListValue) AddFloat(value float32) *ListValue {
	listValue.Value = append(listValue.Value, NewFloatValue(value))
	return listValue
}

// AddString add a string value to slice
func (listValue *ListValue) AddString(value string) *ListValue {
	listValue.Value = append(listValue.Value, NewTextValue(value))
	return listValue
}

// AddString add a string value to slice
func (listValue *ListValue) AddBoolean(value bool) *ListValue {
	listValue.Value = append(listValue.Value, NewBooleanValue(value))
	return listValue
}

// GetString returns string value
func (listValue *ListValue) GetString(i int) string {
	value := listValue.Value[i]
	v, ok := value.(*TextValue)
	if ok {
		return v.ToString()
	}
	return ""
}

// GetFloat returns float value
func (listValue *ListValue) GetFloat(i int) float32 {
	value := listValue.Value[i]
	v, ok := value.(*Float32Value)
	if ok {
		return v.Value
	}
	return 0
}

// GetInt32 returns int32 value
func (listValue *ListValue) GetInt32(i int) int32 {
	value := listValue.Value[i]
	v, ok := value.(*DecimalValue)
	if ok {
		return int32(v.Value)
	}
	return 0
}

// GetInt64 returns int64 value
func (listValue *ListValue) GetInt64(i int) int64 {
	value := listValue.Value[i]
	v, ok := value.(*DecimalValue)
	if ok {
		return v.Value
	}
	return 0
}

// GetBoolean returns bool value
func (listValue *ListValue) GetBoolean(i int) bool {
	value := listValue.Value[i]
	v, ok := value.(*BooleanValue)
	if ok {
		return v.Value
	}
	return false
}

// Size retunrs slice's size
func (listValue *ListValue) Size() int {
	return len(listValue.Value)
}

// Write function writes list value to outputx
func (listValue *ListValue) Write(out *DataOutputX) error {
	valueSize := len(listValue.Value)
	_, err := out.WriteDecimal(int64(valueSize))
	for i := 0; i < valueSize; i++ {
		_, err = out.WriteValue(listValue.Value[i])
	}
	return err
}

// Read function reads value from datainputx
func (listValue *ListValue) Read(in *DataInputX) (Value, error) {
	size, err := in.ReadDecimal()
	for i := int64(0); i < size; i++ {
		var value Value
		value , err = in.ReadValue()
		listValue.Value = append(listValue.Value, value)
	}
	return listValue, err
}

// GetValueType returns value type
func (listValue *ListValue) GetValueType() byte {
	return valueconstants.LIST
}

// ToString returns converted string value
func (listValue *ListValue) ToString() string {
	str := "list value "
	for indx, v := range listValue.Value {
		if v == nil {
			str += strconv.FormatInt(int64(indx), 10) + ": nil "
		} else {
			str += strconv.FormatInt(int64(indx), 10) + ": " + v.ToString() + " "
		}
	}
	return str
}
