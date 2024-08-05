package link

import (
	"sync"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
)

type Link struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}

type lister interface {
	Iterator() (i *iterator.Iterator)                              //创建一个包含链表中所有元素的迭代器并返回其指针
	Sort(Cmp ...comparator.Comparator)                             //将链表中所承载的所有元素进行排序
	Size() (size uint64)                                           //返回链表所承载的元素个数
	Clear()                                                        //清空该链表
	Empty() (b bool)                                               //判断该链表是否位空
	Insert(idx uint64, e interface{})                              //向链表的idx位(下标从0开始)插入元素组e
	Erase(idx uint64)                                              //删除第idx位的元素(下标从0开始)
	Get(idx uint64) (e interface{})                                //获得下标为idx的元素
	Set(idx uint64, e interface{})                                 //在下标为idx的位置上放置元素e
	IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) //返回和元素e相同的第一个下标
	SubList(begin, num uint64) (newList *Link)                     //从begin开始复制最多num个元素以形成新的链表
}

func New() (l *Link) {
	return &Link{
		first: nil,
		last:  nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}

func (l *Link) Iterator() (i *iterator.Iterator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来放入迭代器中
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = iterator.New(&tmp)
	l.mutex.Unlock()
	return i
}

func (l *Link) Sort(Cmp ...comparator.Comparator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来用于排序
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	if len(Cmp) > 0 {
		comparator.Sort(&tmp, Cmp[0])
	} else {
		comparator.Sort(&tmp)
	}
	//将排序结果再放入链表中
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		n.setValue(tmp[idx])
	}
	l.mutex.Unlock()
}

func (l *Link) Size() (size uint64) {
	if l == nil {
		l = New()
	}
	return l.size
}

func (l *Link) Clear() {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//销毁链表
	l.first = nil
	l.last = nil
	l.size = 0
	l.mutex.Unlock()
}

func (l *Link) Empty() (b bool) {
	if l == nil {
		l = New()
	}
	return l.size == 0
}

func (l *Link) Insert(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	n := newNode(e)
	if l.size == 0 {
		//链表中原本无元素,新建链表
		l.first = n
		l.last = n
	} else {
		//链表中存在元素
		if idx == 0 {
			//插入头节点
			n.insertNext(l.first)
			l.first = n
		} else if idx >= l.size {
			//插入尾节点
			l.last.insertNext(n)
			l.last = n
		} else {
			//插入中间节点
			//根据插入的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx-1; i++ {
					m = m.nextNode()
				}
				m.insertNext(n)
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.insertPre(n)
			}
		}
	}
	l.size++
	l.mutex.Unlock()
}

func (l *Link) Erase(idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	if l.size > 0 && idx < l.size {
		//链表中存在元素,且要删除的点在范围内
		if idx == 0 {
			//删除头节点
			l.first = l.first.next
		} else if idx == l.size-1 {
			//删除尾节点
			l.last = l.last.pre
		} else {
			//删除中间节点
			//根据删除的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx; i++ {
					m = m.nextNode()
				}
				m.erase()
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.erase()
			}
		}
		l.size--
		if l.size == 0 {
			//所有节点都被删除,销毁链表
			l.first = nil
			l.last = nil
		}
	}
	l.mutex.Unlock()
}

func (l *Link) Get(idx uint64) (e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return nil
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		e = m.value()
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		e = m.value()
	}
	l.mutex.Unlock()
	return e
}

func (l *Link) Set(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		m.setValue(e)
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		m.setValue(e)
	}
	l.mutex.Unlock()
}

func (l *Link) IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	var equ comparator.Equaler
	if len(Equ) > 0 {
		equ = Equ[0]
	} else {
		equ = comparator.GetEqual()
	}
	n := l.first
	//从头寻找直到找到相等的两个元素即可返回
	for idx = 0; idx < l.size && n != nil; idx++ {
		if equ(n.value(), e) {
			break
		}
		n = n.nextNode()
	}
	l.mutex.Unlock()
	return idx
}

func (l *Link) SubList(begin, num uint64) (newList *Link) {
	if l == nil {
		l = New()
	}
	newList = New()
	l.mutex.Lock()
	if begin < l.size {
		//起点在范围内,可以复制
		n := l.first
		for i := uint64(0); i < begin; i++ {
			n = n.nextNode()
		}
		m := newNode(n.value())
		newList.first = m
		newList.size++
		for i := uint64(0); i < num-1 && i+begin < l.size-1; i++ {
			n = n.nextNode()
			m.insertNext(newNode(n.value()))
			m = m.nextNode()
			newList.size++
		}
		newList.last = m
	}
	l.mutex.Unlock()
	return newList
}
