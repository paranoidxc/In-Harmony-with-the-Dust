package radix

import "strings"

type node struct {
	pattern string
	part    string
	num     int
	sons    map[string]*node
	fuzzy   bool
}

func newNode(part string) *node {
	fuzzy := false
	if len(part) > 0 {
		fuzzy = part[0] == ':' || part[0] == '*'
	}

	return &node{
		pattern: "",
		part:    part,
		num:     0,
		sons:    make(map[string]*node),
		fuzzy:   fuzzy,
	}
}

func analysis(s string) (ss []string, newS string) {
	vs := strings.Split(s, "/")
	ss = make([]string, 0)
	newS = "/"
	for _, item := range vs {
		if item != "" {
			ss = append(ss, item)
			newS = newS + "/" + item
			if item[0] == '*' {
				break
			}
		}
	}
	return ss, newS
}

func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.pattern != "" {
		es = append(es, s+n.part)
	}
	for _, son := range n.sons {
		es = append(es, son.inOrder(s+n.part+"/")...)
	}
	return es
}

func (n *node) insert(pattern string, ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			return false
		}
		n.pattern = pattern
		n.num++
		return true
	}
	s := ss[p]
	son, ok := n.sons[s]
	if !ok {
		son = newNode(s)
		n.sons[s] = son
	}
	b = son.insert(pattern, ss, p+1)
	if b {
		n.num++
	} else {
		if !ok {
			delete(n.sons, s)
		}
	}
	return b
}

func (n *node) erase(ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			n.pattern = ""
			n.num--
			return true
		}
		return false
	}
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return false
	}
	b = son.erase(ss, p+1)
	if b {
		n.num--
		if son.num <= 0 {
			delete(n.sons, s)
		}
	}
	return b
}

func (n *node) delete(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	num = son.delete(ss, p+1)
	if num > 0 {
		son.num -= num
		if son.num <= 0 {
			delete(n.sons, s)
		}
	}
	return num
}

func (n *node) count(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	//从map中找到对应下子结点位置并递归进行查找
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	return son.count(ss, p+1)
}

func (n *node) mate(s string, p int) (m map[string]string, ok bool) {
	searchParts, _ := analysis(s)
	q := n.find(searchParts, 0)
	if q != nil {
		parts, _ := analysis(q.pattern)
		params := make(map[string]string)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return params, true
	}
	return nil, false
}

func (n *node) find(parts []string, height int) (q *node) {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := make([]*node, 0, 0)
	for _, child := range n.sons {
		if child.part == part || child.fuzzy {
			children = append(children, child)
		}
	}
	for _, child := range children {
		result := child.find(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
