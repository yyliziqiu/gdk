package xcq

import (
	"sync"
	"time"
)

type SyncQueue struct {
	qu *Queue
	mu sync.Mutex
}

func NewSync(n int) *SyncQueue {
	return NewSync2(n, "")
}

func NewSync2(n int, path string) *SyncQueue {
	return &SyncQueue{
		qu: New2(n, path),
	}
}

func (t *SyncQueue) EnableDebug() *SyncQueue {
	t.qu.debug = true
	return t
}

func (t *SyncQueue) Get(i int) (any, error) {
	t.mu.Lock()
	item, err := t.qu.Get(i)
	t.mu.Unlock()
	return item, err
}

func (t *SyncQueue) HeadItem() (any, error) {
	t.mu.Lock()
	item, err := t.qu.HeadItem()
	t.mu.Unlock()
	return item, err
}

func (t *SyncQueue) TailItem() (any, error) {
	t.mu.Lock()
	item, err := t.qu.TailItem()
	t.mu.Unlock()
	return item, err
}

func (t *SyncQueue) Status() string {
	t.mu.Lock()
	s := t.qu.Status()
	t.mu.Unlock()
	return s
}

func (t *SyncQueue) Empty() bool {
	t.mu.Lock()
	e := t.qu.Empty()
	t.mu.Unlock()
	return e
}

func (t *SyncQueue) Cap() int {
	t.mu.Lock()
	c := t.qu.Cap()
	t.mu.Unlock()
	return c
}

func (t *SyncQueue) Len() int {
	t.mu.Lock()
	l := t.qu.Len()
	t.mu.Unlock()
	return l
}

func (t *SyncQueue) Push(item any) {
	t.mu.Lock()
	t.qu.Push(item)
	t.mu.Unlock()
}

func (t *SyncQueue) Pop() (any, bool) {
	t.mu.Lock()
	item, ok := t.qu.Pop()
	t.mu.Unlock()
	return item, ok
}

func (t *SyncQueue) Pops(filter Filter) []any {
	t.mu.Lock()
	ret := t.qu.Pops(filter)
	t.mu.Unlock()
	return ret
}

func (t *SyncQueue) Pops2(filter Filter) {
	t.mu.Lock()
	t.qu.Pops2(filter)
	t.mu.Unlock()
}

func (t *SyncQueue) Slide(item any, rmf Remove) (rmd []any) {
	t.mu.Lock()
	rmd = t.qu.Slide(item, rmf)
	t.mu.Unlock()
	return
}

func (t *SyncQueue) SlideN(item any, n int) (rmd []any) {
	t.mu.Lock()
	rmd = t.qu.SlideN(item, n)
	t.mu.Unlock()
	return
}

func (t *SyncQueue) Walk(f func(item any), reverse bool) {
	t.mu.Lock()
	t.qu.Walk(f, reverse)
	t.mu.Unlock()
}

func (t *SyncQueue) Find(filter Filter, reverse bool) (ret any, idx int) {
	t.mu.Lock()
	ret, idx = t.qu.Find(filter, reverse)
	t.mu.Unlock()
	return
}

func (t *SyncQueue) FindAll(f Filter) []any {
	t.mu.Lock()
	ret := t.qu.FindAll(f)
	t.mu.Unlock()
	return ret
}

func (t *SyncQueue) TerminalN(n int, reverse bool) []any {
	t.mu.Lock()
	ret := t.qu.TerminalN(n, reverse)
	t.mu.Unlock()
	return ret
}

func (t *SyncQueue) Terminal(filter Filter, reverse bool) []any {
	t.mu.Lock()
	ret := t.qu.Terminal(filter, reverse)
	t.mu.Unlock()
	return ret
}

func (t *SyncQueue) Window(bgn Filter, end Filter) []any {
	t.mu.Lock()
	ret := t.qu.Window(bgn, end)
	t.mu.Unlock()
	return ret
}

func (t *SyncQueue) Reset(data []any) {
	t.mu.Lock()
	t.qu.Reset(data)
	t.mu.Unlock()
}

func (t *SyncQueue) CopyList() []any {
	t.mu.Lock()
	list := t.qu.CopyList()
	t.mu.Unlock()
	return list
}

func (t *SyncQueue) SaveSnap() error {
	t.mu.Lock()
	err := t.qu.SaveSnap()
	t.mu.Unlock()
	return err
}

func (t *SyncQueue) LoadSnap(item any) error {
	t.mu.Lock()
	err := t.qu.LoadSnap(item)
	t.mu.Unlock()
	return err
}

func (t *SyncQueue) DupSnap(d time.Duration) error {
	t.mu.Lock()
	err := t.qu.DupSnap(d)
	t.mu.Unlock()
	return err
}
