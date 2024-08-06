package cbtree

import "github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"

type node struct {
	value  interface{} //节点中存储的元素
	parent *node       //父节点指针
	left   *node       //左节点指针
	right  *node       //右节点指针
}

func newNode(parent *node, e interface{}) (n *node) {
	return &node{
		value:  e,
		parent: parent,
		left:   nil,
		right:  nil,
	}
}

// node节点做接收者，以前缀序列返回节点集合。
func (n *node) frontOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	es = append(es, n.value)
	if n.left != nil {
		es = append(es, n.left.frontOrder()...)
	}
	if n.right != nil {
		es = append(es, n.right.frontOrder()...)
	}
	return es
}

func (n *node) lastParent(num uint64) (ans *node) {
	if num > 3 {
		//去掉末尾的二进制值
		arr := make([]byte, 0, 64)
		ans = n
		for num > 0 {
			//转化为二进制
			arr = append(arr, byte(num%2))
			num /= 2
		}
		//去掉首位的二进制值
		for i := len(arr) - 2; i >= 1; i-- {
			if arr[i] == 1 {
				ans = ans.right
			} else {
				ans = ans.left
			}
		}
		return ans
	}
	return n
}

func (n *node) insert(num uint64, e interface{}, cmp comparator.Comparator) {
	if n == nil {
		return
	}
	//寻找最后一个父节点
	n = n.lastParent(num)
	//将元素插入最后一个节点
	if num%2 == 0 {
		n.left = newNode(n, e)
		n = n.left
	} else {
		n.right = newNode(n, e)
		n = n.right
	}
	//对插入的节点进行上升
	n.up(cmp)
}

func (n *node) up(cmp comparator.Comparator) {
	if n == nil {
		return
	}
	if n.parent == nil {
		return
	}
	//该节点和父节点都存在
	if cmp(n.parent.value, n.value) > 0 {
		//该节点值小于父节点值,交换两节点值,继续上升
		n.parent.value, n.value = n.value, n.parent.value
		n.parent.up(cmp)
	}
}

func (n *node) delete(num uint64, cmp comparator.Comparator) {
	if n == nil {
		return
	}
	//寻找最后一个父节点
	ln := n.lastParent(num)
	if num%2 == 0 {
		n.value = ln.left.value
		ln.left = nil
	} else {
		n.value = ln.right.value
		ln.right = nil
	}
	//对交换后的节点进行下沉
	n.down(cmp)
}

func (n *node) down(cmp comparator.Comparator) {
	if n == nil {
		return
	}
	if n.right != nil && cmp(n.left.value, n.right.value) >= 0 {
		n.right.value, n.value = n.value, n.right.value
		n.right.down(cmp)
		return
	}
	if n.left != nil && cmp(n.value, n.left.value) >= 0 {
		n.left.value, n.value = n.value, n.left.value
		n.left.down(cmp)
		return
	}
}
