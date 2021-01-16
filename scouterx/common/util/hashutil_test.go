package util

import (
	"testing"
	
	"fmt"
)

func TestStart(t *testing.T) {
	
	fmt.Printf("hash value : %d \n",HashString("hello world"))
	fmt.Printf("hash value : %d \n",HashString("**&^%hello world"))

}