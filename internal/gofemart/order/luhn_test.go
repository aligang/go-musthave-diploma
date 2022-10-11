package order

import (
	"fmt"
	"testing"
)

func TestLuhn(t *testing.T) {
	a := 1000
	for i := 0; i < a; i++ {
		if CheckLuhn(uint64(i)) {
			fmt.Println(i)
		}
	}
}
