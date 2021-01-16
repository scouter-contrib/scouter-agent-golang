package keygen

import (
	"fmt"
	"testing"
)

func TestKeyGen(t *testing.T) {
	fmt.Printf("%d\n", Next())
	fmt.Printf("%d\n", Next())
	fmt.Printf("%d\n", Next())
	fmt.Printf("%d\n", Next())

}
