package stack

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type Stack struct {
	data  []interface{}
	top   uint64 //顶部指针指向实际顶部元素的下一位置
	cap   uint64
	mutex sync.Mutex
}

type stacker interface {
	Iterator() (i *iterator.Iterator) //返回一个包含栈中所有元素的迭代器
	Size() (num uint64)               //返回该栈中元素的使用空间大小
	Clear()                           //清空该栈容器
	Empty() (b bool)                  //判断该栈容器是否为空
	Push(e interface{})               //将元素e添加到栈顶
	Pop()                             //弹出栈顶元素
	Top() (e interface{})             //返回栈顶元素
}

func New() (s *Stack) {
	return &Stack{
		data:  make([]interface{}, 1, 1),
		top:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

func (s *Stack) Iterator() (i *iterator.Iterator) {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	if s.data == nil {
		//data不存在,新建一个
		s.data = make([]interface{}, 1, 1)
		s.top = 0
		s.cap = 1
	} else if s.top < s.cap {
		//释放未使用的空间
		tmp := make([]interface{}, s.top, s.top)
		copy(tmp, s.data)
		s.data = tmp
	}
	//创建迭代器
	i = iterator.New(&s.data)
	s.mutex.Unlock()
	return i
}

func (s *Stack) check() {
	if s == nil {
		s = New()
	}
}

func (s *Stack) Size() (num uint64) {
	s.check()
	return s.top
}

func (s *Stack) Clear() {
	s.check()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = make([]interface{}, 1, 1)
	s.top = 0
	s.cap = 1
}

func (s *Stack) Empty() (b bool) {
	s.check()
	return s.Size() == 0
}

func (s *Stack) Push(e interface{}) {
	s.check()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.top < s.cap {
		//栈未满,直接添加
		s.data[s.top] = e
	} else {
		if s.cap <= 65536 {
			if s.cap == 0 {
				s.cap = 1
			}
			s.cap *= 2
		} else {
			s.cap += 65536
		}
		// 复制
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
		s.data[s.top] = e
	}
	s.top++
}

func (s *Stack) Pop() {
	if s == nil {
		s = New()
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.top--
	if s.cap-s.top >= 65535 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		s.cap -= 65536
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	} else if s.top*2 < s.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		s.cap /= 2
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	}
}

func (s *Stack) Top() (e interface{}) {
	if s == nil {
		return nil
	}
	if s.Empty() {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	e = s.data[s.top-1]
	return e
}
