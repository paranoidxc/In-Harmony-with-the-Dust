package lru

import (
	"container/list"
	"sync"
)

type LRU struct {
	maxBytes int64
	nowBytes int64
	ll       *list.List
	cache    map[string]*list.Element
	onRemove func(key string, value interface{})
	mutex    sync.Mutex
}

type indexes struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

type lruer interface {
	Size() (num int64)                     //返回lru中当前存放的byte数
	Cap() (num int64)                      //返回lru能存放的byte树的最大值
	Clear()                                //清空lru,将其中存储的所有元素都释放
	Empty() (b bool)                       //判断该lru中是否存储了元素
	Insert(key string, value Value)        //向lru中插入以key为索引的value
	Erase(key string)                      //从lru中删除以key为索引的值
	Get(key string) (value Value, ok bool) //从lru中获取以key为索引的value和是否获取成功?
}

func New(maxBytes int64, onRemove func(key string, value interface{})) *LRU {
	return &LRU{
		maxBytes: maxBytes,
		nowBytes: 0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		onRemove: onRemove,
		mutex:    sync.Mutex{},
	}
}

func (l *LRU) Size() (num int64) {
	if l == nil {
		return 0
	}
	return l.nowBytes
}

func (l *LRU) Cap() (num int64) {
	if l == nil {
		return 0
	}
	return l.maxBytes
}

func (l *LRU) Clear() {
	if l == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.cache = make(map[string]*list.Element)
	l.nowBytes = 0
}

func (l *LRU) Empty() (b bool) {
	if l == nil {
		return true
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.nowBytes <= 0
}

func (l *LRU) Insert(key string, value Value) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if ele, ok := l.cache[key]; ok {
		// 该key已存在,直接替换
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		l.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 该key不存在,需要进行插入
		ele := l.ll.PushFront(&indexes{key, value})
		l.cache[key] = ele
		l.nowBytes += int64(len(key)) + int64(value.Len())
	}

	for l.maxBytes != 0 && l.maxBytes < l.nowBytes {
		ele := l.ll.Back()
		if ele != nil {
			l.ll.Remove(ele)
			kv := ele.Value.(*indexes)
			delete(l.cache, kv.key)
			l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
			if l.onRemove != nil {
				l.onRemove(kv.key, kv.value)
			}
		}
	}
}

func (l *LRU) Erase(key string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if ele, ok := l.cache[key]; ok {
		l.ll.Remove(ele)
		kv := ele.Value.(*indexes)
		delete(l.cache, kv.key)
		l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if l.onRemove != nil {
			l.onRemove(kv.key, kv.value)
		}
	}
}

func (l *LRU) Get(key string) (value Value, ok bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		return kv.value, true
	}
	return nil, false
}
