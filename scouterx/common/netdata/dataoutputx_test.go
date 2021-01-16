package netdata

import (
	"fmt"
	"testing"
)

func TestDataOutputX(t *testing.T) {
	out := NewDataOutputX(nil)
	out.WriteInt8(100)
	out.WriteInt8(1)
	out.WriteInt16(13245)
	out.WriteInt32(20022222)
	out.WriteInt64(500033333333)
	out.WriteString("test string....")
	out.WriteString("테스트입니다.")
	out.WriteFloat32(12.456)
	out.WriteDecimal(24000000)
	out.WriteDecimal(35698)
	out.WriteBoolean(true)

	in := NewDataInputX(out.Bytes())
	v1, _ := in.ReadInt8(); fmt.Printf("byte value: %d \n", v1)
	v2, _ := in.ReadInt8(); fmt.Printf("int8 value: %d \n", v2)
	v3, _ := in.ReadInt16(); fmt.Printf("int16 value: %d \n", v3)
	v4, _ := in.ReadInt32(); fmt.Printf("int16 value: %d \n", v4)
	v5, _ := in.ReadInt64(); fmt.Printf("int64 value: %d \n", v5)
	v6, _ := in.ReadString(); fmt.Printf("string value: %s \n", v6)
	v7, _ := in.ReadString(); fmt.Printf("string value: %s \n", v7)
	v8, _ := in.ReadFloat32(); fmt.Printf("float value : %f \n", v8)
	v9, _ := in.ReadDecimal(); fmt.Printf("number value : %d \n", v9)
	v10, _ := in.ReadDecimal(); fmt.Printf("number value : %d \n", v10)
	v11, _ := in.ReadBoolean(); fmt.Printf("bool value : %t \n", v11)

}

func displayPair(a, b interface{}) string {
	return fmt.Sprint(a, b)
}
