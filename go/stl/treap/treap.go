package treap

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"math/rand"
	"sync"
	"time"
)

type Treap struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	rand    *rand.Rand            //随机数生成器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

type treaper interface {
	Iterator() (i *iterator.Iterator) //返回包含该树堆的所有元素,重复则返回多个
	Size() (num int)                  //返回该树堆中保存的元素个数
	Clear()                           //清空该树堆
	Empty() (b bool)                  //判断该树堆是否为空
	Insert(e interface{})             //向树堆中插入元素e
	Erase(e interface{})              //从树堆中删除元素e
	Count(e interface{}) (num int)    //从树堆中寻找元素e并返回其个数
}

func New(isMulti bool, Cmp ...comparator.Comparator) (t *Treap) {
	//设置默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	//创建随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Treap{
		root:    nil,
		size:    0,
		cmp:     cmp,
		rand:    r,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}

func (t *Treap) Iterator() (i *iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	es := t.root.inOrder()
	i = iterator.New(&es)
	return i
}

func (t *Treap) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}

func (t *Treap) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.root = nil
	t.size = 0
}

func (t *Treap) Empty() (b bool) {
	if t == nil {
		return true
	}
	if t.size > 0 {
		return false
	}
	return true
}

func (t *Treap) Insert(e interface{}) {
	//判断容器是否存在
	if t == nil {
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.Empty() {
		//判断比较器是否存在
		if t.cmp == nil {
			t.cmp = comparator.GetCmp(e)
		}
		if t.cmp == nil {
			return
		}
		//插入到根节点
		t.root = newNode(e, t.rand)
		t.size = 1
		return
	}
	//从根节点向下插入
	if t.root.insert(newNode(e, t.rand), t.isMulti, t.cmp) {
		t.size++
	}
}

func (t *Treap) Erase(e interface{}) {
	if t == nil {
		return
	}
	if t.Empty() {
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.size == 1 && t.cmp(t.root.value, e) == 0 {
		//该树堆仅持有一个元素且根节点等价于待删除元素,则将根节点置为nil
		t.root = nil
		t.size = 0
		return
	}
	//从根节点开始删除元素
	if t.root.delete(e, t.isMulti, t.cmp) {
		//删除成功
		t.size--
	}
}

func (t *Treap) Count(e interface{}) (num int) {
	if t == nil {
		//树堆不存在,直接返回0
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	num = t.root.search(e, t.cmp)
	t.mutex.Unlock()
	//树堆存在,从根节点开始查找该元素
	return num
}
