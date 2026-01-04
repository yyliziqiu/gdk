package xcq

import (
	"fmt"
	"testing"
)

var q1 = &Queue{
	step:  10,
	path:  "",
	list:  []any{nil, 1, 2, 3, 4, 5, 6, 7, nil, nil, nil},
	head:  1,
	tail:  8,
	debug: true,
}

var q2 = &Queue{
	step:  10,
	path:  "/private/ws/self/slib/data/q2",
	list:  []any{3, nil, nil, nil, nil, nil, nil, nil, nil, 1, 2},
	head:  9,
	tail:  1,
	debug: true,
}

var q3 = &Queue{
	step:  10,
	path:  "/private/ws/self/slib/data/q3",
	list:  []any{3, 4, 5, 6, 7, 8, nil, nil, nil, 1, 2},
	head:  9,
	tail:  6,
	debug: true,
}

func print1(a ...any) {
	fmt.Println(a...)
}

func print2(a ...any) {
	fmt.Print(a...)
}

func TestLen(t *testing.T) {
	print1("Q1 len: ", q1.len()) // 7
	print1("Q2 len: ", q2.len()) // 3
}

func TestPushAndPop(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
	for _, i := range items {
		q1.Push(i)
	}
	print1(q1.list) // [1 2 3 4 5 6 7 8 9 10 11 12 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]

	for !q1.empty() {
		q1.pop()
	}
	print1(q1.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
}

func TestPushAndPop2(t *testing.T) {
	items := []int{4, 5, 6}
	for _, i := range items {
		q2.Push(i)
	}
	print1(q2.list) // [3 4 5 6 <nil> <nil> <nil> <nil> <nil> 1 2]

	for !q2.empty() {
		q2.pop()
	}
	print1(q2.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]

	items = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.Push(i)
	}
	print1(q2.list) // [8 9 <nil> <nil> 1 2 3 4 5 6 7]

	items = []int{10, 11, 12, 13}
	for _, i := range items {
		q2.Push(i)
	}
	print1(q2.list) // [1 2 3 4 5 6 7 8 9 10 11 12 13 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
}

func TestGet(t *testing.T) {
	print1(q1.Get(0))  // err
	print1(q1.Get(1))  // 1
	print1(q1.Get(4))  // 4
	print1(q1.Get(8))  // err
	print1(q1.Get(10)) // err
}

func TestHeadItem(t *testing.T) {
	print1(q1.HeadItem()) // 1
	print1(q2.HeadItem()) // 1
}

func TestTailItem(t *testing.T) {
	print1(q1.TailItem()) // 7
	print1(q2.TailItem()) // 3
}

func TestEmpty(t *testing.T) {
	print1(q2.Empty()) // false
	for !q2.empty() {
		q2.pop()
	}
	print1(q2.list)    // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	print1(q2.Empty()) // true
}

func TestPops(t *testing.T) {
	result := q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	print1(q1.list) // [<nil> <nil> <nil> <nil> <nil> 5 6 7 <nil> <nil> <nil>]
	print1(result)  // [1 2 3 4]

	result = q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	print1(q1.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	print1(result)  // [5 6 7]

	print1("\n=========================================\n")

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	print1(q3.list) // [<nil> <nil> 5 6 7 8 <nil> <nil> <nil> <nil> <nil>]
	print1(result)  // [1 2 3 4]

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	print1(q3.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	print1(result)  // [5 6 7 8]
}

func TestSlideN(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
	for _, i := range items {
		rm := q1.SlideN(i, 3)
		print1("removed: ", rm)
	}
	print1(q1.list) // [11 12 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> 10]

	items = []int{4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.SlideN(i, 5)
	}
	print1(q2.list) // [<nil> <nil> 5 6 7 8 9 <nil> <nil> <nil> <nil>]
}

func TestSlide(t *testing.T) {
	rm := q1.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	print1("removed: ", rm) // [1 2 3 4]
	print1(q1.list)         // [<nil> <nil> <nil> <nil> <nil> 5 6 7 8 <nil> <nil>]

	rm = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	print1("removed: ", rm) // [1 2 3]
	print1(q2.list)         // [<nil> 8 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]

	rm = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	print1("removed: ", rm) // []
	print1(q2.list)         // [<nil> 8 8 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
}

func TestWalk(t *testing.T) {
	q1.Walk(func(item any) {
		print2(item.(int), " ") // 1 2 3 4 5 6 7
	}, false)
	print1()

	q1.Walk(func(item any) {
		print2(item.(int), " ") // 7 6 5 4 3 2 1
	}, true)
	print1()

	q3.Walk(func(item any) {
		print2(item.(int), " ") // 1 2 3 4 5 6 7 8
	}, false)
	print1()

	q3.Walk(func(item any) {
		print2(item.(int), " ") // 8 7 6 5 4 3 2 1
	}, true)
	print1()
}

func TestFind(t *testing.T) {
	item, _ := q1.Find(func(item any) bool {
		n := item.(int)
		print2(n, " ")
		return n == 3
	}, false)
	print1()
	print1(item) // 3

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		print2(n, " ")
		return n == 100
	}, false)
	print1()
	print1(item) // <nil>

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		print2(n, " ")
		return n == 3
	}, true)
	print1()
	print1(item) // 3

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		print2(n, " ")
		return n == 100
	}, true)
	print1()
	print1(item) // <nil>
}

func TestFindAll(t *testing.T) {
	result := q1.FindAll(func(item any) bool {
		return item.(int) < 5
	})
	print1(result) // [1 2 3 4]
}

func TestTerminalN(t *testing.T) {
	result := q1.TerminalN(3, false)
	print1(result) // [1 2 3]
	result = q1.TerminalN(3, true)
	print1(result) // [7 6 5]

	result = q2.TerminalN(5, false)
	print1(result) // [1 2 3]
	result = q2.TerminalN(5, true)
	print1(result) // [3 2 1]
}

func TestTerminal(t *testing.T) {
	result := q1.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, false)
	print1(result) // [1 2 3]
	result = q1.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, true)
	print1(result) // []

	result = q3.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, false)
	print1(result) // [1 2 3]
	result = q3.Terminal(func(item any) bool {
		return item.(int) >= 5
	}, true)
	print1(result) // [8 7 6 5]
}

func TestWindow(t *testing.T) {
	result := q1.Window(func(item any) bool {
		return item.(int) == 2
	}, func(item any) bool {
		return item.(int) == 4
	})
	print1(result) // [2 3]

	result = q1.Window(func(item any) bool {
		return item.(int) == 20
	}, func(item any) bool {
		return item.(int) == 4
	})
	print1(result) // []

	result = q3.Window(func(item any) bool {
		return item.(int) == 2
	}, func(item any) bool {
		return item.(int) >= 4
	})
	print1(result) // [2 3]

	result = q3.Window(func(item any) bool {
		return item.(int) >= 2
	}, func(item any) bool {
		return item.(int) == 20
	})
	print1(result) // [2 3 4 5 6 7 8]
}

func TestSaveAndLoad(t *testing.T) {
	_ = q3.SaveSnap()
	_ = q3.LoadSnap(1)
	print1(q3.list)
}
