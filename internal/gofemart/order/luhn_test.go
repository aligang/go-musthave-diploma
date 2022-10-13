package order

import (
	"fmt"
	"github.com/theplant/luhn"
	"testing"
)

func TestLuhn(t *testing.T) {
	a := 1000
	for i := 0; i < a; i++ {
		if luhn.Valid(i) {
			fmt.Println(i)
		}
	}
}
