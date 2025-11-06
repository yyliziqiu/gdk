package xutil

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrNoValidValue = errors.New("no valid value")
)

// Round 平滑加权轮询器，在原平滑加权算法基础上优化了权重都相同的情况
type Round struct {
	list []*RoundNode // 节点列表
	len  int          // 节点数量
	swrr bool         // 是否使用平滑加权轮询，如果所有节点的权重全部相同，则循环遍历
	sum  int          // 所有节点的权重不同时，记录所有节点的权重总和
	seq  int32        // 所有节点的权重相同时，计算本次轮询的节点位置
	mu   sync.Mutex   //
}

type RoundNode struct {
	value  any
	weight int
	status int
}

type RoundValue interface {
	GetWeight() int
}

// NewRound 创建轮询器
func NewRound() *Round {
	return &Round{
		list: make([]*RoundNode, 0, 4),
	}
}

// AddValue 添加一个节点
func (t *Round) AddValue(v RoundValue) {
	t.Add(v, v.GetWeight())
}

// Add 添加一个节点
func (t *Round) Add(value any, weight int) {
	t.mu.Lock()
	t.add(value, weight)
	t.mu.Unlock()
}

func (t *Round) add(value any, weight int) {
	// 创建节点
	node := &RoundNode{
		value:  value,
		weight: weight,
	}

	// 判断节点权重是否相同，若权重不完全相同则使用平滑加权，否则使用循环遍历的方式
	if t.len > 1 && node.weight != t.list[t.len-1].weight {
		t.swrr = true
	}

	// 添加节点
	t.list = append(t.list, node)
	t.len = len(t.list)
	t.sum += weight
}

// Next 轮询
func (t *Round) Next() any {
	v, _ := t.NextOrFail()
	return v
}

// NextOrFail 轮询
func (t *Round) NextOrFail() (v any, err error) {
	if t.swrr {
		t.mu.Lock()
		v, err = t.bySwrr()
		t.mu.Unlock()
	} else {
		v, err = t.byLoop()
	}
	return
}

// 平滑加权轮询
func (t *Round) bySwrr() (any, error) {
	var target *RoundNode

	// 将所有节点的状态值加上该节点权重，并选出状态值最大的节点
	for _, node := range t.list {
		node.status += node.weight
		if target == nil || node.status > target.status {
			target = node
		}
	}

	// 无有效节点
	if target == nil {
		return nil, ErrNoValidValue
	}

	// 将选中节点的状态值减去所有节点的权重总和
	target.status -= t.sum

	return target.value, nil
}

// 循环遍历
func (t *Round) byLoop() (any, error) {
	if t.len == 0 {
		return nil, ErrNoValidValue
	}

	i := atomic.AddInt32(&t.seq, 1) & 0x7FFFFFFF

	return t.list[int(i)%t.len].value, nil
}
