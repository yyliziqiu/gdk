package xcq

import (
	"sync"
	"time"
)

type GsQueue struct {
	qu *Queue
	mu sync.Mutex
}

func NewGsQueue(n int) *GsQueue {
	return NewGsQueue2(n, "")
}

func NewGsQueue2(n int, path string) *GsQueue {
	return &GsQueue{
		qu: New2(n, path),
	}
}

func (t *GsQueue) EnableDebug() *GsQueue {
	t.qu.debug = true
	return t
}

func (t *GsQueue) Get(i int) (any, error) {
	t.mu.Lock()
	item, err := t.qu.Get(i)
	t.mu.Unlock()
	return item, err
}

func (t *GsQueue) HeadItem() (any, error) {
	t.mu.Lock()
	item, err := t.qu.HeadItem()
	t.mu.Unlock()
	return item, err
}

func (t *GsQueue) TailItem() (any, error) {
	t.mu.Lock()
	item, err := t.qu.TailItem()
	t.mu.Unlock()
	return item, err
}

func (t *GsQueue) Status() string {
	t.mu.Lock()
	s := t.qu.Status()
	t.mu.Unlock()
	return s
}

func (t *GsQueue) Empty() bool {
	t.mu.Lock()
	e := t.qu.Empty()
	t.mu.Unlock()
	return e
}

func (t *GsQueue) Cap() int {
	t.mu.Lock()
	c := t.qu.Cap()
	t.mu.Unlock()
	return c
}

func (t *GsQueue) Len() int {
	t.mu.Lock()
	l := t.qu.Len()
	t.mu.Unlock()
	return l
}

func (t *GsQueue) Push(item any) {
	t.mu.Lock()
	t.qu.Push(item)
	t.mu.Unlock()
}

func (t *GsQueue) Pop() (any, bool) {
	t.mu.Lock()
	item, ok := t.qu.Pop()
	t.mu.Unlock()
	return item, ok
}

func (t *GsQueue) Pops(filter Filter) []any {
	t.mu.Lock()
	ret := t.qu.Pops(filter)
	t.mu.Unlock()
	return ret
}

func (t *GsQueue) Pops2(filter Filter) {
	t.mu.Lock()
	t.qu.Pops2(filter)
	t.mu.Unlock()
}

func (t *GsQueue) SlideN(item any, n int) (last any, slide bool) {
	t.mu.Lock()
	last, slide = t.qu.SlideN(item, n)
	t.mu.Unlock()
	return
}

func (t *GsQueue) Slide(item any, remove Remove) (last any, n int) {
	t.mu.Lock()
	last, n = t.qu.Slide(item, remove)
	t.mu.Unlock()
	return
}

func (t *GsQueue) Walk(f func(item any), reverse bool) {
	t.mu.Lock()
	t.qu.Walk(f, reverse)
	t.mu.Unlock()
}

func (t *GsQueue) Find(filter Filter, reverse bool) (ret any, idx int) {
	t.mu.Lock()
	ret, idx = t.qu.Find(filter, reverse)
	t.mu.Unlock()
	return
}

func (t *GsQueue) FindAll(f Filter) []any {
	t.mu.Lock()
	ret := t.qu.FindAll(f)
	t.mu.Unlock()
	return ret
}

func (t *GsQueue) TerminalN(n int, reverse bool) []any {
	t.mu.Lock()
	ret := t.qu.TerminalN(n, reverse)
	t.mu.Unlock()
	return ret
}

func (t *GsQueue) Terminal(filter Filter, reverse bool) []any {
	t.mu.Lock()
	ret := t.qu.Terminal(filter, reverse)
	t.mu.Unlock()
	return ret
}

func (t *GsQueue) Window(bgn Filter, end Filter) []any {
	t.mu.Lock()
	ret := t.qu.Window(bgn, end)
	t.mu.Unlock()
	return ret
}

func (t *GsQueue) Reset(data []any) {
	t.mu.Lock()
	t.qu.Reset(data)
	t.mu.Unlock()
}

func (t *GsQueue) CopyList() []any {
	t.mu.Lock()
	list := t.qu.CopyList()
	t.mu.Unlock()

	return list
}

func (t *GsQueue) SaveSnap() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.qu.SaveSnap()
}

func (t *GsQueue) LoadSnap(item any) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.qu.LoadSnap(item)
}

func (t *GsQueue) DupSnap(d time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.qu.DupSnap(d)
}
