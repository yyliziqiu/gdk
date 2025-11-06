package xutil

import (
	"fmt"
	"testing"
)

func TestSwrr(t *testing.T) {
	swrr := NewSwrr2(map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	})

	for i := 0; i < 30; i++ {
		fmt.Print(swrr.Next(), " ")
	}

	fmt.Println()
}
