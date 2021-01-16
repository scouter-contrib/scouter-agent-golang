package netdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapValue1(t *testing.T) {
	ivalue := int64(123)
	tvalue := "test"
	bvalue := true
	mapValue := NewMapValue()
	mapValue.Put("value1", ivalue)
	mapValue.Put("value2", tvalue)
	mapValue.Put("value3", bvalue)

	assert.Equal(t, ivalue, mapValue.GetInt64("value1"))
	assert.Equal(t, tvalue, mapValue.GetString("value2"))
	assert.Equal(t, bvalue, mapValue.GetBoolean("value3"))
	fmt.Printf("%s", mapValue.ToString())

}

func TestMapValue2(t *testing.T) {
	ivalue := int64(123)
	tvalue := "test"
	bvalue := true
	mapValue := NewMapValue()
	mapValue.Put("value1", ivalue)
	mapValue.Put("value2", tvalue)
	mapValue.Put("value3", bvalue)

	out := NewDataOutputX(nil)
	mapValue.Write(out)

	in := NewDataInputX(out.Bytes())

	mapValue2 := NewMapValue()
	mapValue2.Read(in)
	assert.Equal(t, ivalue, mapValue2.GetInt64("value1"))
	assert.Equal(t, tvalue, mapValue2.GetString("value2"))
	assert.Equal(t, bvalue, mapValue2.GetBoolean("value3"))
}
