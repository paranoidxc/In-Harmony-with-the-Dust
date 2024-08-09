package hash_map

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/algo"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/array"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/avl_tree"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type HashMap struct {
	arr   *array.Array //第一层的vector
	hash  algo.Hasher  //hash函数
	size  uint64       //当前存储数量
	cap   uint64       //vector的容量
	mutex sync.Mutex   //并发控制锁
}

type indexes struct {
	key   interface{}
	value interface{}
}

type hashMaper interface {
	Iterator() (i *iterator.Iterator)        //返回一个包含hashMap容器中所有value的迭代器
	Size() (num uint64)                      //返回hashMap已存储的元素数量
	Cap() (num uint64)                       //返回hashMap中的存放空间的容量
	Clear()                                  //清空hashMap
	Empty() (b bool)                         //返回hashMap是否为空
	Insert(key, value interface{}) (b bool)  //向hashMap插入以key为索引的value,若存在会覆盖
	Erase(key interface{}) (b bool)          //删除hashMap中以key为索引的value
	GetKeys() (keys []interface{})           //返回hashMap中所有的keys
	Get(key interface{}) (value interface{}) //以key为索引寻找vlue
}

func New(hash ...algo.Hasher) (hm *HashMap) {
	var h algo.Hasher
	if len(hash) == 0 {
		h = nil
	} else {
		h = hash[0]
	}
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	//新建vector并将其扩容到16
	v := array.New()
	for i := 0; i < 16; i++ {
		//vector中嵌套avl树
		v.PushBack(avl_tree.New(false, cmp))
	}
	return &HashMap{
		arr:   v,
		hash:  h,
		size:  0,
		cap:   16,
		mutex: sync.Mutex{},
	}
}

func (hm *HashMap) Iterator() (i *iterator.Iterator) {
	if hm == nil {
		return nil
	}
	if hm.arr == nil {
		return nil
	}
	hm.mutex.Lock()
	//取出hashMap中存放的所有value
	values := make([]interface{}, 0, 1)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avl_tree.AvlTree)
		ite := avl.Iterator()
		es := make([]interface{}, 0, 1)
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idx := j.Value().(*indexes)
			es = append(es, idx.value)
		}
		values = append(values, es...)
	}
	//将所有value放入迭代器中
	i = iterator.New(&values)
	hm.mutex.Unlock()
	return i
}

func (hm *HashMap) Size() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.size
}

func (hm *HashMap) Cap() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.cap
}

func (hm *HashMap) Clear() {
	if hm == nil {
		return
	}
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	//重建vector并扩容到16
	v := array.New()
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	for i := 0; i < 16; i++ {
		v.PushBack(avl_tree.New(false, cmp))
	}
	hm.arr = v
	hm.size = 0
	hm.cap = 16
}

func (hm *HashMap) Empty() (b bool) {
	if hm == nil {
		return false
	}
	return hm.size > 0
}

func (hm *HashMap) Insert(key, value interface{}) (b bool) {
	if hm == nil {
		return false
	}
	if hm.arr == nil {
		return false
	}
	if hm.hash == nil {
		hm.hash = algo.GetHash(key)
	}
	if hm.hash == nil {
		return false
	}

	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	//计算hash值并找到对应的avl树
	hash := hm.hash(key) % hm.cap
	avl := hm.arr.At(hash).(*avl_tree.AvlTree)
	idx := &indexes{
		key:   key,
		value: value,
	}
	//判断是否存在该avl树中
	if avl.Count(idx) == 0 {
		//avl树中不存在相同key,插入即可
		avl.Insert(idx)
		hm.size++
		if hm.size >= hm.cap/4*3 {
			//当达到扩容条件时候进行扩容
			hm.expend()
		}
	} else {
		//覆盖
		avl.Insert(idx)
	}
	return true
}

func (hm *HashMap) expend() {
	//取出所有的key-value
	idxs := make([]*indexes, 0, hm.size)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avl_tree.AvlTree)
		ite := avl.Iterator()
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idxs = append(idxs, j.Value().(*indexes))
		}
	}
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	//对vector进行扩容,扩容到其容量上限即可
	hm.arr.PushBack(avl_tree.New(false, cmp))
	for i := uint64(0); i < hm.arr.Size()-1; i++ {
		hm.arr.At(i).(*avl_tree.AvlTree).Clear()
	}
	for i := hm.arr.Size(); i < hm.arr.Cap(); i++ {
		hm.arr.PushBack(avl_tree.New(false, cmp))
	}
	//将vector容量设为hashMap容量
	hm.cap = hm.arr.Cap()
	//重新将所有的key-value插入到hashMap中去
	for i := 0; i < len(idxs); i++ {
		key, value := idxs[i].key, idxs[i].value
		hash := hm.hash(key) % hm.cap
		avl := hm.arr.At(hash).(*avl_tree.AvlTree)
		idx := &indexes{
			key:   key,
			value: value,
		}
		avl.Insert(idx)
	}
}

func (hm *HashMap) Erase(key interface{}) (b bool) {
	if hm == nil {
		return false
	}
	if hm.arr == nil {
		return false
	}
	if hm.hash == nil {
		return false
	}
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	//计算该key的hash值
	hash := hm.hash(key) % hm.cap
	avl := hm.arr.At(hash).(*avl_tree.AvlTree)
	idx := &indexes{
		key:   key,
		value: nil,
	}
	//从对应的avl树中删除该key-value
	b = avl.Erase(idx)
	if b {
		//删除成功,此时size-1,同时进行缩容判断
		hm.size--
		if hm.size < hm.cap/8*3 && hm.cap > 16 {
			hm.shrink()
		}
	}
	return b
}

func (hm *HashMap) shrink() {
	//取出所有key-value
	idxs := make([]*indexes, 0, hm.size)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avl_tree.AvlTree)
		ite := avl.Iterator()
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idxs = append(idxs, j.Value().(*indexes))
		}
	}
	//进行缩容,当vector的cap与初始不同时,说明缩容结束
	cap := hm.arr.Cap()
	for cap == hm.arr.Cap() {
		hm.arr.PopBack()
	}
	hm.cap = hm.arr.Cap()
	//将所有的key-value重新放入hashMap中
	for i := 0; i < len(idxs); i++ {
		key, value := idxs[i].key, idxs[i].value
		hash := hm.hash(key) % hm.cap
		avl := hm.arr.At(hash).(*avl_tree.AvlTree)
		idx := &indexes{
			key:   key,
			value: value,
		}
		avl.Insert(idx)
	}
}

func (hm *HashMap) GetKeys() (keys []interface{}) {
	if hm == nil {
		return nil
	}
	if hm.arr == nil {
		return nil
	}
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	keys = make([]interface{}, 0, 1)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avl_tree.AvlTree)
		ite := avl.Iterator()
		es := make([]interface{}, 0, 1)
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idx := j.Value().(*indexes)
			es = append(es, idx.key)
		}
		keys = append(keys, es...)
	}
	return keys
}

func (hm *HashMap) Get(key interface{}) (value interface{}) {
	if hm == nil {
		return
	}
	if hm.arr == nil {
		return
	}
	if hm.hash == nil {
		hm.hash = algo.GetHash(key)
	}
	if hm.hash == nil {
		return
	}
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	//计算hash值
	hash := hm.hash(key) % hm.cap
	//从avl树中找到对应该hash值的key-value
	info := hm.arr.At(hash).(*avl_tree.AvlTree).Find(&indexes{key: key, value: nil})
	if info == nil {
		return nil
	}
	return info.(*indexes).value
}
