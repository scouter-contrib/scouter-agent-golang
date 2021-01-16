package netdata

import (
	"fmt"
	"testing"
)

func TestTextValue(t *testing.T) {
	value := NewTextValue("aaaa")

	fmt.Printf("value: %s \n", value.Value)

}
