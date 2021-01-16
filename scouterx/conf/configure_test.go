package conf

import (
	"fmt"
	"github.com/magiconair/properties"
	"testing"
)

func TestConfigLoad(t *testing.T)  {
	conf := GetInstance()
	fmt.Println(conf.NetCollectorIP)
}

func TestProperties(t *testing.T)  {
	var s = `
a=1
b=2=3
c=4
empty=
e=hollo,:ok=312/323/<pp>
empty_no=
b1=true
b2=false
b3=
`
	p := properties.MustLoadString(s)
	fmt.Println(p.MustGetInt("a"))
	fmt.Println(p.MustGetString("b"))
	fmt.Println(p.MustGetInt("c"))
	strEmpty := p.GetString("empty", "default")
	fmt.Println("strEmpty=" + strEmpty)
	strNoKey := p.GetString("nokey", "default")
	fmt.Println("strNoKey=" + strNoKey)
	fmt.Println(p.GetString("e", "default"))
	fmt.Println(p.GetInt("empty_no", 10))

	fmt.Println(p.GetBool("b1", false))
	fmt.Println(p.GetBool("b2", true))
	fmt.Println(p.GetBool("b3", true))
}
