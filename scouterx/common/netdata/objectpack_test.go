package netdata

import (
	"fmt"
	"testing"
)

func TestObjectPack(t *testing.T) {
	objPack := NewObjectPack2()
	objPack.SiteID = "abdec"
	objPack.ObjHash = 12345
	objPack.ObjName = "testObjName"
	objPack.ObjType = "container"
	objPack.Address = "1.1.1.1"
	objPack.Family = 1
	objPack.Version = "v1"
	objPack.Wakeup = 1234565698
	objPack.Tags.Put("key1", "test")
	fmt.Printf("%s\n", objPack.ToString())

	out := NewDataOutputX(nil)
	objPack.Write(out)
	in := NewDataInputX(out.Bytes())
	objPack2 := NewObjectPack2()
	objPack2.Read(in)
	fmt.Printf("%s\n", objPack2.ToString())

}
