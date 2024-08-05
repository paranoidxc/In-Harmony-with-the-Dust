package ring

import (
	"sync"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
)

type Ring struct {
	current *node
	size    uint64

	mutex sync.Mutex
}

type ringer interface {
	Iterator() (i *iterator.Iterator) //创建一个包含环中所有元素的迭代器并返回其指针
	Size() (size uint64)              //返回环所承载的元素个数
	Clear()                           //清空该环
	Empty() (b bool)                  //判断该环是否位空
	Insert(e interface{})             //向环当前位置后方插入元素e
	Erase()                           //删除当前结点并持有下一结点
	Value() (e interface{})           //返回当前持有结点的元素
	Set(e interface{})                //在当前结点设置其承载的元素为e
	Next()                            //持有下一节点
	Pre()                             //持有上一结点
}

func New() (r *Ring) {
	return &Ring{
		current: nil,
		size:    0,
		mutex:   sync.Mutex{},
	}
}

func check(r *Ring) *Ring {
	if r == nil {
		return New()
	}
	return r
}

func (r *Ring) Iterator() (i *iterator.Iterator) {
	if r == nil {
		r = New()
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	tmp := make([]interface{}, r.size, r.size)
	for n, idx := r.current, uint64(0); n != nil && idx < r.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = iterator.New(&tmp)
	return i
}

func (r *Ring) Size() (size uint64) {
	if r == nil {
		r = New()
	}

	return r.size
}

func (r *Ring) Clear() {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.current = nil
	r.size = 0
}

func (r *Ring) Empty() (b bool) {
	if r == nil {
		r = New()
	}
	return r.size == 0
}

func (r *Ring) Insert(e interface{}) {
	r = check(r)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	n := newNode(e)
	if r.size == 0 {
		r.current = n
	} else {
		r.current.insertNext(n)
	}

	r.size += 1
}

func (r *Ring) Erase() {
	r = check(r)
	if r.size == 0 {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.size == 1 {
		r.current = nil
	} else {
		//环内还有其他结点,将持有结点后移一位
		//后移后将当前结点前插原持有结点的前结点
		nextNode := r.current.nextNode()
		prevNode := r.current.preNode()
		prevNode.next = nextNode
		nextNode.pre = prevNode
		r.current = nextNode
	}

	r.size -= 1
}

func (r *Ring) Value() (e interface{}) {
	r = check(r)
	if r.current == nil {
		return nil
	}
	return r.current.value()
}

func (r *Ring) Set(e interface{}) {
	r = check(r)
	if r.current == nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.current.setValue(e)
}

func (r *Ring) Next() {
	r = check(r)
	if r.current == nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.current = r.current.nextNode()
}

func (r *Ring) Pre() {
	r = check(r)
	if r.current == nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.current = r.current.preNode()
}
