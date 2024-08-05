package queue

import (
	"sync"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
)

type Queue struct {
	data  []interface{} //泛型切片
	begin uint64        //首节点下标
	end   uint64        //尾节点下标
	cap   uint64        //容量
	mutex sync.Mutex    //并发控制锁
}

type queuer interface {
	Iterator() (i *iterator.Iterator) //返回包含队列中所有元素的迭代器
	Size() (num uint64)               //返回该队列中元素的使用空间大小
	Clear()                           //清空该队列
	Empty() (b bool)                  //判断该队列是否为空
	Push(e interface{})               //将元素e添加到该队列末尾
	Pop() (e interface{})             //将该队列首元素弹出并返回
	Front() (e interface{})           //获取该队列首元素
	Back() (e interface{})            //获取该队列尾元素
}

func New() (q *Queue) {
	return &Queue{
		data:  make([]interface{}, 1, 1),
		begin: 0,
		end:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

func (q *Queue) Iterator() (i *iterator.Iterator) {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	tmp := make([]interface{}, q.end-q.begin, q.end-q.begin)
	copy(tmp, q.data[q.begin:q.end])
	i = iterator.New(&tmp)
	q.mutex.Unlock()
	return i
}

func (q *Queue) Size() (num uint64) {
	if q == nil {
		q = New()
	}
	return q.end - q.begin
}

func (q *Queue) Clear() {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	q.data = make([]interface{}, 1, 1)
	q.begin = 0
	q.end = 0
	q.cap = 1
	q.mutex.Unlock()
}

func (q *Queue) Empty() (b bool) {
	if q == nil {
		q = New()
	}
	return q.Size() <= 0
}

func (q *Queue) Push(e interface{}) {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	if q.end < q.cap {
		//不需要扩容
		q.data[q.end] = e
	} else {
		//需要扩容
		if q.begin > 0 {
			//首部有冗余,整体前移
			for i := uint64(0); i < q.end-q.begin; i++ {
				q.data[i] = q.data[i+q.begin]
			}
			q.end -= q.begin
			q.begin = 0
		} else {
			//冗余不足,需要扩容
			if q.cap <= 65536 {
				//容量翻倍
				if q.cap == 0 {
					q.cap = 1
				}
				q.cap *= 2
			} else {
				//容量增加2^16
				q.cap += 2 ^ 16
			}
			//复制扩容前的元素
			tmp := make([]interface{}, q.cap, q.cap)
			copy(tmp, q.data)
			q.data = tmp
		}
		q.data[q.end] = e
	}
	q.end++
	q.mutex.Unlock()
}

func (q *Queue) Pop() (e interface{}) {
	if q == nil {
		q = New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.begin]
	q.begin++
	if q.begin >= 1024 || q.begin*2 > q.end {
		//首部冗余超过2^10或首部冗余超过实际使用
		q.cap -= q.begin
		q.end -= q.begin
		tmp := make([]interface{}, q.cap, q.cap)
		copy(tmp, q.data[q.begin:])
		q.data = tmp
		q.begin = 0
	}
	q.mutex.Unlock()
	return e
}

func (q *Queue) Front() (e interface{}) {
	if q == nil {
		q = New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.begin]
}

func (q *Queue) Back() (e interface{}) {
	if q == nil {
		q = New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.end-1]
}
