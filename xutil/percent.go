package xutil

import (
	"sync"
)

// Percent 按设定比例返回 true
type Percent struct {
	cycle int        // 此值越大 reset 周期越大
	bound float64    // Next() 返回 true 的比例
	total int        // Next() 调用次数
	count int        // Next() 返回 true 的次数
	mu    sync.Mutex //
}

func NewPercent(p float64) *Percent {
	return &Percent{
		cycle: 1e6,
		bound: p / 100,
		total: 0,
		count: 0,
	}
}

func (t *Percent) Next() bool {
	t.mu.Lock()
	ret := t.do()
	t.mu.Unlock()

	return ret
}

func (t *Percent) do() bool {
	if t.total > t.cycle {
		t.total = 0
		t.count = 0
	}

	t.total++
	if float64(t.count)/float64(t.total) >= t.bound {
		return false
	}

	t.count++

	return true
}
