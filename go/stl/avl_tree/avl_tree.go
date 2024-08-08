package avl_tree

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type AvlTree struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

type avlTreer interface {
	Iterator() (i *iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num int)                  //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{})             //向二叉树中插入元素e
	Erase(e interface{})              //从二叉树中删除元素e
	Count(e interface{}) (num int)    //从二叉树中寻找元素e并返回其个数
}

func New(isMulti bool, cmps ...comparator.Comparator) (avl *AvlTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	return &AvlTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
	}
}

func (avl *AvlTree) Iterator() (i *iterator.Iterator) {
	if avl == nil {
		return nil
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()

	es := avl.root.inOrder()
	i = iterator.New(&es)
	return i
}

func (avl *AvlTree) Size() (num int) {
	if avl == nil {
		return 0
	}
	return avl.size
}

func (avl *AvlTree) Clear() {
	if avl == nil {
		return
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()

	avl.root = nil
	avl.size = 0
}

func (avl *AvlTree) Empty() (b bool) {
	if avl == nil {
		return true
	}
	if avl.size > 0 {
		return false
	}
	return true
}

func (n *node) leftRotate() (m *node) {
	//左旋转
	headNode := n.right
	n.right = headNode.left
	headNode.left = n
	//更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

func (n *node) rightRotate() (m *node) {
	//右旋转
	headNode := n.left
	n.left = headNode.right
	headNode.right = n
	//更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

func (n *node) rightLeftRotate() (m *node) {
	//右旋转,之后左旋转
	//以失衡点右结点先右旋转
	sonHeadNode := n.right.rightRotate()
	n.right = sonHeadNode
	//再以失衡点左旋转
	return n.leftRotate()
}

func (n *node) leftRightRotate() (m *node) {
	//左旋转,之后右旋转
	//以失衡点左结点先左旋转
	sonHeadNode := n.left.leftRotate()
	n.left = sonHeadNode
	//再以失衡点左旋转
	return n.rightRotate()
}

func (avl *AvlTree) Insert(e interface{}) {
	if avl == nil {
		return
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()

	if avl.Empty() {
		if avl.cmp == nil {
			avl.cmp = comparator.GetCmp(e)
		}
		if avl.cmp == nil {
			return
		}
		//二叉树为空,用根节点承载元素e
		avl.root = newNode(e)
		avl.size = 1
		return
	}
	//从根节点进行插入,并返回节点,同时返回是否插入成功
	var b bool
	avl.root, b = avl.root.insert(e, avl.isMulti, avl.cmp)
	if b {
		//插入成功,数量+1
		avl.size++
	}
}

func (avl *AvlTree) Erase(e interface{}) {
	if avl == nil {
		return
	}
	if avl.Empty() {
		return
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()

	if avl.size == 1 && avl.cmp(avl.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		avl.root = nil
		avl.size = 0
		return
	}
	//从根节点进行插入,并返回节点,同时返回是否删除成功
	var b bool
	avl.root, b = avl.root.erase(e, avl.cmp)
	if b {
		avl.size--
	}
}

func (avl *AvlTree) Count(e interface{}) (num int) {
	if avl == nil {
		//二叉树为空,返回0
		return 0
	}
	if avl.Empty() {
		return 0
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()
	num = avl.root.count(e, avl.isMulti, avl.cmp)
	return num
}

func (avl *AvlTree) Find(e interface{}) (ans interface{}) {
	if avl == nil {
		//二叉树为空,返回0
		return 0
	}
	if avl.Empty() {
		return 0
	}
	avl.mutex.Lock()
	defer avl.mutex.Unlock()
	ans = avl.root.find(e, avl.isMulti, avl.cmp)
	return ans
}
