package treap

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"math/rand"
)

type node struct {
	value    interface{} //节点中存储的元素
	priority uint32      //该节点的优先级,随机生成
	num      int         //该节点中存储的数量
	left     *node       //左节点指针
	right    *node       //右节点指针
}

func newNode(e interface{}, rand *rand.Rand) (n *node) {
	return &node{
		value:    e,
		priority: uint32(rand.Intn(4294967295)),
		num:      1,
		left:     nil,
		right:    nil,
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

func (n *node) rightRotate() {
	if n == nil {
		return
	}
	if n.left == nil {
		return
	}
	//新建节点作为更换后的n节点
	tmp := &node{
		value:    n.value,
		priority: n.priority,
		num:      n.num,
		left:     n.left.right,
		right:    n.right,
	}
	//原n节点左节点上移到n节点位置
	n.right = tmp
	n.value = n.left.value
	n.priority = n.left.priority
	n.num = n.left.num
	n.left = n.left.left
}

func (n *node) leftRotate() {
	if n == nil {
		return
	}
	if n.right == nil {
		return
	}
	//新建节点作为更换后的n节点
	tmp := &node{
		value:    n.value,
		priority: n.priority,
		num:      n.num,
		left:     n.left,
		right:    n.right.left,
	}
	//原n节点右节点上移到n节点位置
	n.left = tmp
	n.value = n.right.value
	n.priority = n.right.priority
	n.num = n.right.num
	n.right = n.right.right
}

func (n *node) insert(e *node, isMulti bool, cmp comparator.Comparator) (b bool) {
	if cmp(n.value, e.value) > 0 {
		if n.left == nil {
			//将左节点直接设为e
			n.left = e
			b = true
		} else {
			//对左节点进行递归插入
			b = n.left.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行右转
			n.rightRotate()
		}
		return b
	} else if cmp(n.value, e.value) < 0 {
		if n.right == nil {
			//将右节点直接设为e
			n.right = e
			b = true
		} else {
			//对右节点进行递归插入
			b = n.right.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行左转
			n.leftRotate()
		}
		return b
	}
	if isMulti {
		//允许重复
		n.num++
		return true
	}
	//不允许重复,对值进行覆盖
	n.value = e.value
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
			//待删除节点无子节点,直接删除即可
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
			//待删除节点无子节点,直接删除即可
			if n.left.left == nil && n.left.right == nil {
				//左子树可直接删除
				n.left = nil
				return true
			}
		}
		//从左子树继续删除
		return n.left.delete(e, isMulti, cmp)
	}
	if isMulti && n.num > 1 {
		//允许重复且数量超过1
		n.num--
		return true
	}
	//删除该节点
	tmp := n
	//左右子节点都存在则选择优先级较小一个进行旋转
	for tmp.left != nil && tmp.right != nil {
		if tmp.left.priority < tmp.right.priority {
			tmp.rightRotate()
			if tmp.right.left == nil && tmp.right.right == nil {
				tmp.right = nil
				return false
			}
			tmp = tmp.right
		} else {
			tmp.leftRotate()
			if tmp.left.left == nil && tmp.left.right == nil {
				tmp.left = nil
				return false
			}
			tmp = tmp.left
		}
	}
	if tmp.left == nil && tmp.right != nil {
		//到左子树为nil时直接换为右子树即可
		tmp.value = tmp.right.value
		tmp.num = tmp.right.num
		tmp.priority = tmp.right.priority
		tmp.left = tmp.right.left
		tmp.right = tmp.right.right
	} else if tmp.right == nil && tmp.left != nil {
		//到右子树为nil时直接换为左子树即可
		tmp.value = tmp.left.value
		tmp.num = tmp.left.num
		tmp.priority = tmp.left.priority
		tmp.right = tmp.left.right
		tmp.left = tmp.left.left
	}
	//当左右子树都为nil时直接结束
	return true
}

func (n *node) search(e interface{}, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	if cmp(n.value, e) > 0 {
		return n.left.search(e, cmp)
	} else if cmp(n.value, e) < 0 {
		return n.right.search(e, cmp)
	}
	return n.num
}
