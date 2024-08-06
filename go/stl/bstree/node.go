package bstree

import "github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"

type node struct {
	value interface{} //节点中存储的元素
	num   uint64      //该元素数量
	left  *node       //左节点指针
	right *node       //右节点指针
}

func newNode(e interface{}) (n *node) {
	return &node{
		value: e,
		num:   1,
		left:  nil,
		right: nil,
	}
}

func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//n中承载元素小于e,从右子树继续插入
	if cmp(n.value, e) < 0 {
		if n.right == nil {
			//右子树为nil,直接插入右子树即可
			n.right = newNode(e)
			return true
		} else {
			return n.right.insert(e, isMulti, cmp)
		}
	}
	//n中承载元素大于e,从左子树继续插入
	if cmp(n.value, e) > 0 {
		if n.left == nil {
			//左子树为nil,直接插入左子树即可
			n.left = newNode(e)
			return true
		} else {
			return n.left.insert(e, isMulti, cmp)
		}
	}
	//n中承载元素等于e
	if isMulti {
		//允许重复
		n.num++
		return true
	}
	//不允许重复,直接进行覆盖
	n.value = e
	return false
}

func (n *node) delete(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//n中承载元素小于e,从右子树继续删除
	if cmp(n.value, e) < 0 {
		if n.right == nil {
			//右子树为nil,删除终止
			return false
		}
		if cmp(e, n.right.value) == 0 && (!isMulti || n.right.num == 1) {
			if n.right.left == nil && n.right.right == nil {
				//右子树可直接删除
				n.right = nil
				return true
			}
		}
		//从右子树继续删除
		return n.right.delete(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续删除
	if cmp(n.value, e) > 0 {
		if n.left == nil {
			//左子树为nil,删除终止
			return false
		}
		if cmp(e, n.left.value) == 0 && (!isMulti || n.left.num == 1) {
			if n.left.left == nil && n.left.right == nil {
				//左子树可直接删除
				n.left = nil
				return true
			}
		}
		//从左子树继续删除
		return n.left.delete(e, isMulti, cmp)
	}
	//n中承载元素等于e
	if (*n).num > 1 && isMulti {
		//允许重复且个数超过1,则减少num即可
		(*n).num--
		return true
	}
	if n.left == nil && n.right == nil {
		//该节点无前缀节点和后继节点,删除即可
		*(&n) = nil
		return true
	}
	if n.left != nil {
		//该节点有前缀节点,找到前缀节点进行交换并删除即可
		ln := n.left
		if ln.right == nil {
			n.value = ln.value
			n.num = ln.num
			n.left = ln.left
		} else {
			for ln.right.right != nil {
				ln = ln.right
			}
			n.value = ln.right.value
			n.num = ln.right.num
			ln.right = ln.right.left
		}
	} else if (*n).right != nil {
		//该节点有后继节点,找到后继节点进行交换并删除即可
		tn := n.right
		if tn.left == nil {
			n.value = tn.value
			n.num = tn.num
			n.right = tn.right
		} else {
			for tn.left.left != nil {
				tn = tn.left
			}
			n.value = tn.left.value
			n.num = tn.left.num
			tn.left = tn.left.right
		}
		return true
	}
	return true
}

func (n *node) search(e interface{}, isMulti bool, cmp comparator.Comparator) (num uint64) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.search(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.search(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}
