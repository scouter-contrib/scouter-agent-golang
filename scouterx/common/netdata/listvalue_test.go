package netdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListValue1(t *testing.T) {
	listValue := NewListValue()
	ivalue := int32(100)
	lvalue := int64(4000000)
	tvalue := "test"
	listValue.AddInt32(ivalue)
	listValue.AddInt64(lvalue)
	listValue.AddString(tvalue)
	assert.Equal(t, ivalue, listValue.GetInt32(0))
	assert.Equal(t, lvalue, listValue.GetInt64(1))
	assert.Equal(t, tvalue, listValue.GetString(2))
	fmt.Printf("list value %s", listValue.ToString())
	/*
		out := NewDataOutputX()
		listValue.Write(out)
		in :=
	*/
}

func TestListValue2(t *testing.T) {
	listValue := NewListValue()
	ivalue := int32(100)
	lvalue := int64(4000000)
	tvalue := "test"
	listValue.AddInt32(ivalue)
	listValue.AddInt64(lvalue)
	listValue.AddString(tvalue)
	out := NewDataOutputX(nil)
	listValue.Write(out)
	in := NewDataInputX(out.Bytes())

	listValue2 := NewListValue()
	listValue2.Read(in)
	assert.Equal(t, ivalue, listValue2.GetInt32(0))
	assert.Equal(t, lvalue, listValue2.GetInt64(1))
	assert.Equal(t, tvalue, listValue2.GetString(2))

}
