package array

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
	"sync"
)

type Array struct {
	data  []interface{} //动态数组
	len   uint64        //当前已用数量
	cap   uint64        //可容纳元素数量
	mutex sync.Mutex    //并发控制锁
}

// vector扩容边界,边界内进行翻倍扩容,边界外进行固定扩容
const bound = 4294967296

type vectorer interface {
	Iterator() (i *iterator.Iterator)  //返回一个包含vector所有元素的迭代器
	Sort(Cmp ...comparator.Comparator) //利用比较器对其进行排序
	Size() (num uint64)                //返回vector的长度
	Cap() (num uint64)                 //返回vector的容量
	Clear()                            //清空vector
	Empty() (b bool)                   //返回vector是否为空,为空则返回true反之返回false
	PushBack(e interface{})            //向vector末尾插入一个元素
	PopBack()                          //弹出vector末尾元素
	Insert(idx uint64, e interface{})  //向vector第idx的位置插入元素e,同时idx后的其他元素向后退一位
	Erase(idx uint64)                  //删除vector的第idx个元素
	Reverse()                          //逆转vector中的数据顺序
	At(idx uint64) (e interface{})     //返回vector的第idx的元素
	Front() (e interface{})            //返回vector的第一个元素
	Back() (e interface{})             //返回vector的最后一个元素
}

func New() (v *Array) {
	return &Array{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

func (v *Array) Iterator() (i *iterator.Iterator) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用的空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
	}
	//创建迭代器
	i = iterator.New(&v.data)
	v.mutex.Unlock()
	return i
}

func (v *Array) Sort(Cmp ...comparator.Comparator) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
		v.cap = v.len
	}
	//调用比较器的Sort进行排序
	if len(Cmp) == 0 {
		comparator.Sort(&v.data)
	} else {
		comparator.Sort(&v.data, Cmp[0])
	}
	v.mutex.Unlock()
}

func (v *Array) Size() (num uint64) {
	if v == nil {
		v = New()
	}
	return v.len
}

func (v *Array) Cap() (num uint64) {
	if v == nil {
		v = New()
	}
	return v.cap
}

func (v *Array) Clear() {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	//清空data
	v.data = make([]interface{}, 1, 1)
	v.len = 0
	v.cap = 1
	v.mutex.Unlock()
}

func (v *Array) Empty() (b bool) {
	if v == nil {
		v = New()
	}
	return v.Size() <= 0
}

func (v *Array) PushBack(e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.len < v.cap {
		//还有冗余,直接添加
		v.data[v.len] = e
	} else {
		//冗余不足,需要扩容
		if v.cap <= bound {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加bound
			v.cap += bound
		}
		//复制扩容前的元素
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
		v.data[v.len] = e
	}
	v.len++
	v.mutex.Unlock()
}

func (v *Array) PopBack() {
	if v == nil {
		v = New()
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	v.len--
	if v.cap-v.len >= bound {
		//容量和实际使用差值超过bound时,容量直接减去bound
		v.cap -= bound
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	} else if v.len*2 < v.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		v.cap /= 2
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	v.mutex.Unlock()
}

func (v *Array) Insert(idx uint64, e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	var tmp []interface{}
	if v.len >= v.cap {
		//冗余不足,进行扩容
		if v.cap <= bound {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加bound
			v.cap += bound
		}
		//复制扩容前的元素
		tmp = make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	//从后往前复制,即将idx后的全部后移一位即可
	var p uint64
	for p = v.len; p > 0 && p > uint64(idx); p-- {
		v.data[p] = v.data[p-1]
	}
	v.data[p] = e
	v.len++
	v.mutex.Unlock()
}

func (v *Array) Erase(idx uint64) {
	if v == nil {
		v = New()
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	for p := idx; p < v.len-1; p++ {
		v.data[p] = v.data[p+1]
	}
	v.len--
	if v.cap-v.len >= bound {
		//容量和实际使用差值超过bound时,容量直接减去bound
		v.cap -= bound
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	} else if v.len*2 < v.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		v.cap /= 2
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	v.mutex.Unlock()
}

func (v *Array) Reverse() {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用的空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
		v.cap = v.len
	}
	for i := uint64(0); i < v.len/2; i++ {
		v.data[i], v.data[v.len-i-1] = v.data[v.len-i-1], v.data[i]
	}
	v.mutex.Unlock()
}

func (v *Array) At(idx uint64) (e interface{}) {
	if v == nil {
		v = New()
		return nil
	}
	v.mutex.Lock()
	if idx < 0 && idx >= v.Size() {
		v.mutex.Unlock()
		return nil
	}
	if v.Size() > 0 {
		e = v.data[idx]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}

func (v *Array) Front() (e interface{}) {
	if v == nil {
		v = New()
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[0]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}

func (v *Array) Back() (e interface{}) {
	if v == nil {
		v = New()
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[v.len-1]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}
