package avl_tree

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
)

type node struct {
	value interface{} //节点中存储的元素
	num   int         //该元素数量
	depth int         //该节点的深度
	left  *node       //左节点指针
	right *node       //右节点指针
}

func newNode(e interface{}) (n *node) {
	return &node{
		value: e,
		num:   1,
		depth: 1,
		left:  nil,
		right: nil,
	}
}

func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := 0; i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}

func (n *node) getDepth() (depth int) {
	if n == nil {
		return 0
	}
	return n.depth
}

func max(a, b int) (m int) {
	if a > b {
		return a
	} else {
		return b
	}
}

func (n *node) getMin() (e interface{}, num int) {
	if n == nil {
		return nil, 0
	}
	if n.left == nil {
		return n.value, n.num
	} else {
		return n.left.getMin()
	}
}

func (n *node) adjust() (m *node) {
	if n.right.getDepth()-n.left.getDepth() >= 2 {
		//右子树高于左子树且高度差超过2,此时应当对n进行左旋
		if n.right.right.getDepth() > n.right.left.getDepth() {
			//由于右右子树高度超过右左子树,故可以直接左旋
			n = n.leftRotate()
		} else {
			//由于右右子树不高度超过右左子树
			//所以应该先右旋右子树使得右子树高度不超过左子树
			//随后n节点左旋
			n = n.rightLeftRotate()
		}
	} else if n.left.getDepth()-n.right.getDepth() >= 2 {
		//左子树高于右子树且高度差超过2,此时应当对n进行右旋
		if n.left.left.getDepth() > n.left.right.getDepth() {
			//由于左左子树高度超过左右子树,故可以直接右旋
			n = n.rightRotate()
		} else {
			//由于左左子树高度不超过左右子树
			//所以应该先左旋左子树使得左子树高度不超过右子树
			//随后n节点右旋
			n = n.leftRightRotate()
		}
	}
	return n
}

func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (m *node, b bool) {
	//节点不存在,应该创建并插入二叉树中
	if n == nil {
		return newNode(e), true
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续插入
		n.left, b = n.left.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	if cmp(e, n.value) > 0 {
		//从右子树继续插入
		n.right, b = n.right.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	//该节点元素与待插入元素相同
	if isMulti {
		//允许重复,数目+1
		n.num++
		return n, true
	}
	//不允许重复,对值进行覆盖
	n.value = e
	return n, false
}

func (n *node) erase(e interface{}, cmp comparator.Comparator) (m *node, b bool) {
	if n == nil {
		//待删除值不存在,删除失败
		return n, false
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续删除
		n.left, b = n.left.erase(e, cmp)
	} else if cmp(e, n.value) > 0 {
		//从右子树继续删除
		n.right, b = n.right.erase(e, cmp)
	} else if cmp(e, n.value) == 0 {
		//存在相同值,从该节点删除
		b = true
		if n.num > 1 {
			//有重复值,节点无需删除,直接-1即可
			n.num--
		} else {
			//该节点需要被删除
			if n.left != nil && n.right != nil {
				//找到该节点后继节点进行交换删除
				n.value, n.num = n.right.getMin()
				//从右节点继续删除,同时可以保证删除的节点必然无左节点
				n.right, b = n.right.erase(n.value, cmp)
			} else if n.left != nil {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	//当n节点仍然存在时,对其进行调整
	if n != nil {
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		n = n.adjust()
	}
	return n, b
}

func (n *node) count(e interface{}, isMulti bool, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.count(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.count(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}

func (n *node) find(e interface{}, isMulti bool, cmp comparator.Comparator) (ans interface{}) {
	if n == nil {
		return nil
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.find(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.find(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.value
}
