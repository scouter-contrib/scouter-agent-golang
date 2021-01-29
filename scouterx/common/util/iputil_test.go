package util

import (
	"testing"
	
	"fmt"
)

func TestIpToBytes(t *testing.T) {

	bytes := IpToBytes("118.241.190.59")
	fmt.Printf("bytes : %s \n",string(bytes))
}
