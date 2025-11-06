package xutil

type Trie struct {
	root *Node
}

type Node struct {
	Data any
	Leaf bool
	Next map[byte]*Node
}

func New() *Trie {
	return &Trie{
		root: &Node{
			Leaf: false,
			Next: map[byte]*Node{},
		},
	}
}

func (t *Trie) BatchAdd(data map[string]any) {
	for k, v := range data {
		t.Add(k, v)
	}
}

func (t *Trie) Add(prefix string, data any) {
	if prefix == "" {
		return
	}

	curr := t.root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		next, ok := curr.Next[c]
		if !ok {
			next = &Node{Next: map[byte]*Node{}}
			curr.Next[c] = next
		}
		curr = next
	}

	curr.Data = data
	curr.Leaf = true
}

// Exist 判断是否存在指定字符串
func (t *Trie) Exist(str string) (any, bool) {
	curr := t.root
	for i := 0; i < len(str); i++ {
		c := str[i]
		next, ok := curr.Next[c]
		if !ok {
			return nil, false
		}
		curr = next
	}

	if !curr.Leaf {
		return nil, false
	}

	return curr.Data, true
}

// Match 判断是否存在指定字符串的前缀，最长匹配
// n 最多匹配位数
func (t *Trie) Match(str string, n int) (any, bool) {
	var data any

	curr := t.root
	for i := 0; i < len(str) && i <= n; i++ {
		c := str[i]
		next, ok := curr.Next[c]
		if !ok {
			break
		}
		if next.Leaf {
			data = next.Data
		}
		curr = next
	}

	return data, data != nil
}
