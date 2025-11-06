package xutil

import (
	"fmt"
	"testing"
)

func TestPercent(t *testing.T) {
	p := NewPercent(34)

	a, c := 0, 5235252
	for i := 0; i < c; i++ {
		ok := p.Next()
		if ok {
			a++
		}
	}

	fmt.Printf("a = %d, c = %f\n", a, float64(a)/float64(c))
}
