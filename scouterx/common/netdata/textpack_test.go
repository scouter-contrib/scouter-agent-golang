package netdata

import (
	"fmt"
	"testing"
)

func TestTextPack(t *testing.T) {
	pack := NewTextPack()
	pack.Xtype = "error"
	pack.Hash = 12345
	pack.Text = "testText"
	fmt.Printf("%s\n", pack.ToString())

	out := NewDataOutputX(nil)
	pack.Write(out)

}
