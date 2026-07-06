package xutil

type Trie struct {
	root *Node
}

type Node struct {
	Data any
	Leaf bool
	Next map[byte]*Node
}

func NewTrie() *Trie {
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

	prev := t.root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		curr, ok := prev.Next[c]
		if !ok {
			curr = &Node{Next: map[byte]*Node{}}
			prev.Next[c] = curr
		}
		prev = curr
	}

	prev.Data = data
	prev.Leaf = true
}

// Exist 判断 Tire 树中是否存在指定字符串的前缀
func (t *Trie) Exist(str string) (any, bool) {
	prev := t.root
	for i := 0; i < len(str); i++ {
		curr, ok := prev.Next[str[i]]
		if !ok {
			return nil, false
		}
		if curr.Leaf {
			return curr.Data, true
		}
		prev = curr
	}
	return nil, false
}

// Match 判断 Tire 树中是否存在指定字符串的前缀，最长匹配
// n 最多匹配位数
func (t *Trie) Match(str string, n int) (any, bool) {
	var data any

	prev := t.root
	for i := 0; i <= n && i < len(str); i++ {
		curr, ok := prev.Next[str[i]]
		if !ok {
			break
		}
		if curr.Leaf {
			data = curr.Data
		}
		prev = curr
	}

	return data, data != nil
}
