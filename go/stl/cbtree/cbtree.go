package cbtree

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type CBTree struct {
	root  *node                 //根节点指针
	size  uint64                //存储元素数量
	cmp   comparator.Comparator //比较器
	mutex sync.Mutex            //并发控制锁
}

type cbTreeer interface {
	Iterator() (i *iterator.Iterator) //返回包含该二叉树的所有元素
	Size() (num uint64)               //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Push(e interface{})               //向二叉树中插入元素e
	Pop()                             //从二叉树中弹出顶部元素
	Top() (e interface{})             //返回该二叉树的顶部元素
}

func New(Cmp ...comparator.Comparator) (cb *CBTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	return &CBTree{
		root:  nil,
		size:  0,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

func (cb *CBTree) Iterator() (i *iterator.Iterator) {
	if cb == nil {
		cb = New()
	}
	cb.mutex.Lock()
	es := cb.root.frontOrder()
	i = iterator.New(&es)
	cb.mutex.Unlock()
	return i
}

func (cb *CBTree) Size() (num uint64) {
	if cb == nil {
		cb = New()
	}
	return cb.size
}

func (cb *CBTree) Clear() {
	if cb == nil {
		cb = New()
	}
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.root = nil
	cb.size = 0
}

func (cb *CBTree) Empty() (b bool) {
	if cb == nil {
		cb = New()
	}
	return cb.size == 0
}

func (cb *CBTree) Push(e interface{}) {
	if cb == nil {
		cb = New()
	}

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.Empty() {
		if cb.cmp == nil {
			cb.cmp = comparator.GetCmp(e)
		}
		if cb.cmp == nil {
			return
		}
		cb.root = newNode(nil, e)
		cb.size++
	} else {
		cb.size++
		cb.root.insert(cb.size, e, cb.cmp)
	}
}

func (cb *CBTree) Pop() {
	if cb == nil {
		return
	}
	if cb.Empty() {
		return
	}

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.size == 1 {
		//该二叉树仅剩根节点,直接删除即可
		cb.root = nil
	} else {
		//该二叉树删除根节点后还有其他节点可生为跟节点
		cb.root.delete(cb.size, cb.cmp)
	}
	cb.size--
}

func (cb *CBTree) Top() (e interface{}) {
	if cb == nil {
		cb = New()
	}

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	e = cb.root.value
	return e
}
