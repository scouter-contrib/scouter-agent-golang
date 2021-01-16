package netdata

import (
	"fmt"
	"testing"
)

func TestPerfCounterPack(t *testing.T) {
	pack := NewPerfCounterPack()
	pack.Put("fvalue", 1.23)
	pack.Put("ivalue", 123)
	pack.Put("tvalue", " test value")
	pack.Put("bvalue", false)
	fmt.Printf("%s", pack.ToString())
}

func TestPerfCounterPack2(t *testing.T) {
	pack := NewPerfCounterPack()
	pack.Put("fvalue", 1.23)
	pack.Put("ivalue", 123)
	pack.Put("tvalue", " test value2")
	pack.Put("bvalue", false)
	out := NewDataOutputX(nil)
	pack.Write(out)
	in := NewDataInputX(out.Bytes())
	pack2 := NewPerfCounterPack()
	pack2.Read(in)
	fmt.Printf("%s", pack2.ToString())
}
