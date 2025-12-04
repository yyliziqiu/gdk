package xutil

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRound(t *testing.T) {
	r := NewRound()
	r.add("A", 10)
	r.add("B", 20)
	r.add("C", 30)
	r.add("D", 40)
	r.add("E", 0)

	fmt.Println("start: ", r.swrr)
	t1 := time.Now().UnixMilli()

	var a, b, c, d, e int32
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10000000; j++ {
				k := r.Next()
				s := k.(string)
				switch s {
				case "A":
					atomic.AddInt32(&a, 1)
				case "B":
					atomic.AddInt32(&b, 1)
				case "C":
					atomic.AddInt32(&c, 1)
				case "D":
					atomic.AddInt32(&d, 1)
				case "E":
					atomic.AddInt32(&e, 1)
				default:
					panic("unreachable")
				}
			}
		}()
	}
	wg.Wait()

	fmt.Println("end, cost: ", time.Now().UnixMilli()-t1)
	fmt.Println("A: ", a, " B: ", b, " C: ", c, " D: ", d, " E: ", e)
}

func TestRound2(t *testing.T) {
	r := NewRound()
	r.Add("A", 100)
	r.Add("B", 0)
	r.Add("C", 0)

	fmt.Println(r.swrr)

	for i := 0; i < 20; i++ {
		k := r.Next()
		fmt.Print(k.(string), " ")
	}

	fmt.Println()
}

func TestRound3(t *testing.T) {
	r := NewRound()
	r.Add("A", 100)

	fmt.Println(r.swrr)

	for i := 0; i < 20; i++ {
		k := r.Next()
		if k != nil {
			fmt.Print(k.(string), " ")
		}
	}

	fmt.Println()
}
