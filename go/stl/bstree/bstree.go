package bstree

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type BSTree struct {
	root    *node                 //根节点指针
	size    uint64                //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

type bSTreeer interface {
	Iterator() (i *iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num uint64)               //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{})             //向二叉树中插入元素e
	Erase(e interface{})              //从二叉树中删除元素e
	Count(e interface{}) (num uint64) //从二叉树中寻找元素e并返回其个数
}

func New(isMulti bool, Cmp ...comparator.Comparator) (bs *BSTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) == 0 {
		cmp = nil
	} else {
		cmp = Cmp[0]
	}
	return &BSTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}

func (bs *BSTree) Iterator() (i *iterator.Iterator) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	es := bs.root.inOrder()
	i = iterator.New(&es)
	bs.mutex.Unlock()
	return i
}

func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := uint64(0); i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}

func (bs *BSTree) Size() (num uint64) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size
}

func (bs *BSTree) Clear() {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.root = nil
	bs.size = 0
}

func (bs *BSTree) Empty() (b bool) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size == 0
}

func (bs *BSTree) Insert(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if bs.Empty() {
		//二叉树为空,用根节点承载元素e
		if bs.cmp == nil {
			bs.cmp = comparator.GetCmp(e)
		}
		if bs.cmp == nil {
			return
		}
		bs.root = newNode(e)
		bs.size++
		return
	}
	//二叉树不为空,从根节点开始查找添加元素e
	if bs.root.insert(e, bs.isMulti, bs.cmp) {
		bs.size++
	}
}

func (bs *BSTree) Erase(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	if bs.size == 0 {
		return
	}
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	if bs.size == 1 && bs.cmp(bs.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		bs.root = nil
		bs.size = 0
		return
	}
	//从根节点开始删除元素e
	//如果删除成功则将size-1
	if bs.root.delete(e, bs.isMulti, bs.cmp) {
		bs.size--
	}
}

func (bs *BSTree) Count(e interface{}) (num uint64) {
	if bs == nil {
		//二叉树不存在,返回0
		return 0
	}
	if bs.Empty() {
		//二叉树为空,返回0
		return 0
	}
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	//从根节点开始查找并返回查找结果
	num = bs.root.search(e, bs.isMulti, bs.cmp)
	return num
}
