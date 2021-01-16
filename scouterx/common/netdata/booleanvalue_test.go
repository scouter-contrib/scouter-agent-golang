package netdata

import (
	"testing"
)

func TestBooleanVale(t *testing.T) {
	out := NewDataOutputX(nil)
	bvalue := NewBooleanValue(true)
	bvalue.Write(out)
	in := NewDataInputX(out.Bytes())
	testValue, _ := in.ReadBoolean()
	if testValue != true {
		t.Error("test error")
	}

}
