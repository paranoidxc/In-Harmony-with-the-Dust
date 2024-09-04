package radix

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type radix struct {
	root  *node
	size  int
	mutex sync.Mutex
}

type radixer interface {
	Iterator() (i *iterator.Iterator)
	Size() (num int)
	Clear()
	Empty() (b bool)
	Insert(s string) (b bool)
	Erase(s string) (b bool)
	Delete(s string) (num int)
	Count(s string) (num int)
	Mate(s string) (m map[string]string, ok bool)
}

func New() (r *radix) {
	return &radix{root: newNode(""),
		size:  0,
		mutex: sync.Mutex{},
	}
}

func (r *radix) Iterator() (i *iterator.Iterator) {
	if r == nil {
		return nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	es := r.root.inOrder("")
	i = iterator.New(&es)
	return i
}

func (r *radix) Size() (num int) {
	if r == nil {
		return 0
	}
	if r.root == nil {
		return 0
	}
	return r.size
}

func (r *radix) Clear() {
	if r == nil {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.root = newNode("")
	r.size = 0
}

func (r *radix) Empty() (b bool) {
	if r == nil {
		return true
	}
	return r.size == 0
}

func (r *radix) Insert(s string) (b bool) {
	if r == nil {
		return false
	}
	ss, s := analysis(s)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.root == nil {
		r.root = newNode("")
	}

	b = r.root.insert(s, ss, 0)
	if b {
		//插入成功,size+1
		r.size++
	}
	return b
}

func (r *radix) Erase(s string) (b bool) {
	if r.Empty() {
		return false
	}
	if len(s) == 0 {
		return false
	}
	if r.root == nil {
		return false
	}

	ss, _ := analysis(s)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	b = r.root.erase(ss, 0)
	if b {
		//删除成功,size-1
		r.size--
		if r.size == 0 {
			//所有string都被删除,根节点置为nil
			r.root = nil
		}
	}
	return b
}

func (r *radix) Delete(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	if r.root == nil {
		return 0
	}
	ss, _ := analysis(s)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	num = r.root.delete(ss, 0)
	if num > 0 {
		r.size -= num
		if r.size <= 0 {
			r.root = nil
		}
	}
	return num
}

func (r *radix) Count(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if r.root == nil {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	num = r.root.count(ss, 0)
	r.mutex.Unlock()
	return num
}

func (r *radix) Mate(s string) (m map[string]string, ok bool) {
	if r.Empty() {
		return nil, false
	}
	if len(s) == 0 {
		return nil, false
	}
	if r.root == nil {
		return nil, false
	}
	m, ok = r.root.mate(s, 0)
	return m, ok
}
