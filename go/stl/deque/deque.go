package deque

import (
	"sync"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
)

type Deque struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}

type dequer interface {
	Iterator() (i *iterator.Iterator) //返回包含双向队列中所有元素的迭代器
	Size() (size uint64)              //返回该双向队列中元素的使用空间大小
	Clear()                           //清空该双向队列
	Empty() (b bool)                  //判断该双向队列是否为空
	PushFront(e interface{})          //将元素e添加到该双向队列的首部
	PushBack(e interface{})           //将元素e添加到该双向队列的尾部
	PopFront()                        //将该双向队列首元素弹出
	PopBack()                         //将该双向队列首元素弹出
	Front() (e interface{})           //获取该双向队列首部元素
	Back() (e interface{})            //获取该双向队列尾部元素
}

func New() *Deque {
	return &Deque{
		first: nil,
		last:  nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}

func (d *Deque) Iterator() (i *iterator.Iterator) {
	if d == nil {
		d = New()
	}
	tmp := make([]interface{}, 0, d.size)
	//遍历链表的所有节点,将其中承载的元素全部复制出来
	for m := d.first; m != nil; m = m.nextNode() {
		tmp = append(tmp, m.value()...)
	}
	return iterator.New(&tmp)
}

func (d *Deque) Size() (size uint64) {
	if d == nil {
		d = New()
	}
	return d.size
}

func (d *Deque) Clear() {
	if d == nil {
		d = New()
		return
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.first = nil
	d.last = nil
	d.size = 0
}

func (d *Deque) Empty() (b bool) {
	if d == nil {
		d = New()
	}

	return d.Size() == 0
}

func (d *Deque) PushFront(e interface{}) {
	if d == nil {
		d = New()
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.size++

	if d.first == nil {
		d.first = createFirst()
		d.last = d.first
	}

	d.first = d.first.pushFront(e)
}

func (d *Deque) PushBack(e interface{}) {
	if d == nil {
		d = New()
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.size++

	if d.last == nil {
		d.last = createLast()
		d.first = d.last
	}

	d.last = d.last.pushBack(e)
}

func (d *Deque) PopFront() (e interface{}) {
	if d == nil {
		d = New()
	}

	if d.size == 0 {
		return
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	//利用首节点删除首元素
	//返回新的首节点
	//fmt.Println("d.first", d.first)
	//fmt.Println("d.last", d.last)
	e = d.first.front()
	d.first = d.first.popFront()
	//fmt.Println("e:", e)

	d.size--
	if d.size == 0 {
		//全部删除完成,释放空间,并将首尾节点设为nil
		d.first = nil
		d.last = nil
	}

	return e
}

func (d *Deque) PopBack() (e interface{}) {
	if d == nil {
		d = New()
	}
	if d.size == 0 {
		return nil
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	//利用尾节点删除首元素
	//返回新的尾节点
	e = d.last.back()
	d.last = d.last.popBack()
	d.size--
	if d.size == 0 {
		d.first = nil
		d.last = nil
	}
	return e
}

func (d *Deque) Front() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.first.front()
}

func (d *Deque) Back() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.last.back()
}
