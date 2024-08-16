package trie

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type Trie struct {
	root  *node      //根节点指针
	size  int        //存放的元素数量
	mutex sync.Mutex //并发控制锁
}

type trieer interface {
	Iterator() (i *iterator.Iterator)        //返回包含该trie的所有string
	Size() (num int)                         //返回该trie中保存的元素个数
	Clear()                                  //清空该trie
	Empty() (b bool)                         //判断该trie是否为空
	Insert(s string, e interface{}) (b bool) //向trie中插入string并携带元素e
	Erase(s string) (b bool)                 //从trie中删除以s为索引的元素e
	Delete(s string) (num int)               //从trie中删除以s为前缀的所有元素
	Count(s string) (num int)                //从trie中寻找以s为前缀的string单词数
	Find(s string) (e interface{})           //从trie中寻找以s为索引的元素e
}

func New() (t *Trie) {
	return &Trie{
		root:  newNode(nil),
		size:  0,
		mutex: sync.Mutex{},
	}
}

func (t *Trie) Iterator() (i *iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	//找到trie中存在的所有string
	es := t.root.inOrder("")
	i = iterator.New(&es)
	t.mutex.Unlock()
	return i
}

func (t *Trie) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}

func (t *Trie) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = newNode(nil)
	t.size = 0
	t.mutex.Unlock()
}

func (t *Trie) Empty() (b bool) {
	if t == nil {
		return true
	}
	return t.size == 0
}

func (t *Trie) Insert(s string, e interface{}) (b bool) {
	if t == nil {
		return
	}
	if len(s) == 0 {
		return false
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.root == nil {
		//避免根节点为nil
		t.root = newNode(nil)
	}
	//从根节点开始插入
	b = t.root.insert(s, 0, e)
	if b {
		//插入成功,size+1
		t.size++
	}
	return b
}

func (t *Trie) Erase(s string) (b bool) {
	if t == nil {
		return false
	}
	if t.Empty() {
		return false
	}
	if len(s) == 0 {
		//长度为0无法删除
		return false
	}
	if t.root == nil {
		//根节点为nil即无法删除
		return false
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//从根节点开始删除
	b = t.root.erase(s, 0)
	if b {
		//删除成功,size-1
		t.size--
		if t.size == 0 {
			//所有string都被删除,根节点置为nil
			t.root = nil
		}
	}
	return b
}

func (t *Trie) Delete(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return 0
	}
	if len(s) == 0 {
		//长度为0无法删除
		return 0
	}
	if t.root == nil {
		//根节点为nil即无法删除
		return 0
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//从根节点开始删除
	num = t.root.delete(s, 0)
	if num > 0 {
		//删除成功
		t.size -= num
		if t.size <= 0 {
			//所有string都被删除,根节点置为nil
			t.root = nil
		}
	}
	return num
}

func (t *Trie) Count(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return 0
	}
	if t.root == nil {
		return 0
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//统计所有以s为前缀的string的数量并返回
	num = int(t.root.count(s, 0))
	return num
}

func (t *Trie) Find(s string) (e interface{}) {
	if t == nil {
		return nil
	}
	if t.Empty() {
		return nil
	}
	if t.root == nil {
		return nil
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	//从根节点开始查找以s为索引的元素e
	e = t.root.find(s, 0)
	return e
}
