package util

import (
	"fmt"
	"testing"
)

func TestCopyArray(t *testing.T) {
	src := []byte{10, 20, 30, 40, 50, 60}
	target := CopyArray(src, 2, 4)
	for _, v := range target {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}
