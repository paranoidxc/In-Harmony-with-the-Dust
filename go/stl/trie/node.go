package trie

type node struct {
	num      int         //以当前结点为前缀的string的数量
	children [64]*node   //分叉
	value    interface{} //当前结点承载的元素
}

func newNode(e interface{}) (n *node) {
	return &node{
		num:   0,
		value: e,
	}
}

func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es = append(es, s)
	}
	for i, p := 0, 0; i < 64 && p < n.num; i++ {
		if n.children[i] != nil {
			if i < 26 {
				es = append(es, n.children[i].inOrder(s+string(i+'a'))...)
			} else if i < 52 {
				es = append(es, n.children[i].inOrder(s+string(i-26+'A'))...)
			} else if i == 62 {
				es = append(es, n.children[i].inOrder(s+string('+'))...)
			} else if i == 63 {
				es = append(es, n.children[i].inOrder(s+string('/'))...)
			} else {
				es = append(es, n.children[i].inOrder(s+string(i-52+'0'))...)
			}
			p++
		}
	}
	return es
}

func getIdx(c byte) (idx int) {
	if c >= 'a' && c <= 'z' {
		idx = int(c - 'a')
	} else if c >= 'A' && c <= 'Z' {
		idx = int(c-'A') + 26
	} else if c >= '0' && c <= '9' {
		idx = int(c-'0') + 52
	} else if c == '+' {
		idx = 62
	} else if c == '/' {
		idx = 63
	} else {
		idx = -1
	}
	return idx
}

func (n *node) insert(s string, p int, e interface{}) (b bool) {
	if p == len(s) {
		if n.value != nil {
			return false
		}
		n.value = e
		n.num++
		return true
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.children[idx] == nil {
		n.children[idx] = newNode(nil)
	}
	b = n.children[idx].insert(s, p+1, e)
	if b {
		n.num++
	}
	return b
}

func (n *node) erase(s string, p int) (b bool) {
	if p == len(s) {
		if n.value != nil {
			n.value = nil
			n.num--
			return true
		}
		return false
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.children[idx] == nil {
		return false
	}
	b = n.children[idx].erase(s, p+1)
	if b {
		n.num--
		if n.children[idx].num == 0 {
			n.children[idx] = nil
		}
	}
	return b
}

func (n *node) delete(s string, p int) (num int) {
	if p == len(s) {
		return n.num
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return 0
	}
	if n.children[idx] == nil {
		return 0
	}
	num = n.children[idx].delete(s, p+1)
	if num > 0 {
		n.num -= num
		if n.children[idx].num <= 0 {
			n.children[idx] = nil
		}
	}
	return num
}

func (n *node) count(s string, p int) (num int) {
	if p == len(s) {
		return n.num
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return 0
	}
	if n.children[idx] == nil {
		return 0
	}
	return n.children[idx].count(s, p+1)
}

func (n *node) find(s string, p int) (e interface{}) {
	if p == len(s) {
		return n.value
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return nil
	}
	if n.children[idx] == nil {
		return nil
	}
	return n.children[idx].find(s, p+1)
}
