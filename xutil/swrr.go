package xutil

import (
	"sync"
)

// 算法流程如下：
// 1. 所有节点的当前权重初始值均为零：{0,0,0}
// 2. 开始轮询，所有节点的当前权重值加上该节点的设定权重值，并选取更新后当前权重值最大的节点作为命中节点
// 3. 将命中节点的当前权重值减去所有节点设定权重的总和作为其新权重值，并将命中节点返回
// 4. 下次轮询重复第2步骤
//
// 设 A、B、C 三个节点的权重分别为：4、2、1，演算步骤如下：
// 步骤	选择前当前值	选择节点	选择后当前值
// 1	{ 4, 2, 1}	A	    {-3, 2, 1}
// 2	{ 1, 4, 2}	B	    { 1,-3, 2}
// 3	{ 5,-1, 3}	A	    {-2,-1, 3}
// 4	{ 2, 1, 4}	C	    { 2, 1,-3}
// 5	{ 6, 3,-2}	A	    {-1, 3,-2}
// 6	{ 3, 5,-1}	B	    { 3,-2,-1}
// 7	{ 7, 0, 0}	A	    { 0, 0, 0}
//
// 三个节点的命中次数符合 4:2:1，而且权重大的节点不会霸占选择权
// 经过一个周期(七轮选择)后，当前权重值又回到了{0, 0, 0}
// 以上过程将按照周期进行循环，完全符合我们先前期望的平滑性

// Swrr Smooth Weight Round Robin
type Swrr[T any] struct {
	nodes []*SwrrNode[T]
	total int
	mu    sync.Mutex
}

type SwrrNode[T any] struct {
	value  T
	weight int // 设定权重
	status int // 当前权重
}

// NewSwrr 创建平滑权重轮询器
func NewSwrr[T any]() *Swrr[T] {
	return &Swrr[T]{}
}

// NewSwrr2 创建平滑权重轮询器
func NewSwrr2[T comparable](weights map[T]int) *Swrr[T] {
	t := &Swrr[T]{
		nodes: make([]*SwrrNode[T], 0, len(weights)),
	}
	for k, v := range weights {
		t.Add(k, v)
	}
	return t
}

// Add 添加一个权重节点
func (t *Swrr[T]) Add(value T, weight int) {
	t.mu.Lock()
	t.nodes = append(t.nodes, &SwrrNode[T]{value: value, weight: weight})
	t.total += weight
	t.mu.Unlock()
}

// Next 加权轮询
func (t *Swrr[T]) Next() T {
	t.mu.Lock()
	v := t.next()
	t.mu.Unlock()

	return v
}

func (t *Swrr[T]) next() T {
	var best *SwrrNode[T]

	// 将所有节点的当前权重值加上设定权重，并选出当前权重值最大的节点
	for _, node := range t.nodes {
		node.status += node.weight
		if best == nil || node.status > best.status {
			best = node
		}
	}

	if best == nil {
		return SwrrNode[T]{}.value
	}

	// 将选中节点的当前权重值减去设定权重值总和
	best.status -= t.total

	return best.value
}
