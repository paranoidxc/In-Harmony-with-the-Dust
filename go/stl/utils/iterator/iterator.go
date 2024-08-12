package iterator

import "fmt"

type Iterator struct {
	data  *[]interface{} //该迭代器中存放的元素集合的指针
	index int            //该迭代器当前指向的元素下标，-1即不存在元素
}

type Iteratorer interface {
	Begin() (I *Iterator)      //将该迭代器设为位于首节点并返回新迭代器
	End() (I *Iterator)        //将该迭代器设为位于尾节点并返回新迭代器
	Get(idx int) (I *Iterator) //将该迭代器设为位于第idx节点并返回该迭代器
	Value() (e interface{})    //返回该迭代器下标所指元素
	HasNext() (b bool)         //判断该迭代器是否可以后移
	Next() (b bool)            //将该迭代器后移一位
	HasPre() (b bool)          //判罚该迭代器是否可以前移
	Pre() (b bool)             //将该迭代器前移一位
}

func New(data *[]interface{}, Idx ...int) (i *Iterator) {
	//迭代器下标
	var idx int
	if len(Idx) <= 0 {
		//没有传入下标，则将下标设为0
		idx = 0
	} else {
		//有传入下标，则将传入下标第一个设为迭代器下标
		idx = Idx[0]
	}
	if len((*data)) > 0 {
		//如果元素集合非空，则判断下标是否超过元素集合范围
		if idx >= len((*data)) {
			//如果传入下标超过元素集合范围则寻找最近的下标值
			idx = len((*data)) - 1
		}
	} else {
		//如果元素集合为空则将下标设为-1
		idx = -1
	}
	//新建并返回迭代器
	return &Iterator{
		data:  data,
		index: idx,
	}
}

func (i *Iterator) Begin() (I *Iterator) {
	if i == nil {
		//迭代器为空，直接结束
		return nil
	}
	if len((*i.data)) == 0 {
		//迭代器元素集合为空，下标设为-1
		i.index = -1
	} else {
		//迭代器元素集合非空，下标设为0
		i.index = 0
	}
	//返回修改后的新指针
	return &Iterator{
		data:  i.data,
		index: i.index,
	}
}

func (i *Iterator) End() (I *Iterator) {
	if i == nil {
		//迭代器为空，直接返回
		return nil
	}
	if len((*i.data)) == 0 {
		//元素集合为空，下标设为-1
		i.index = -1
	} else {
		//元素集合非空，下标设为最后一个元素的下标
		i.index = len((*i.data)) - 1
	}
	//返回修改后的该指针
	return &Iterator{
		data:  i.data,
		index: i.index,
	}
}

func (i *Iterator) Get(idx int) (I *Iterator) {
	if i == nil {
		//迭代器为空，直接返回
		return nil
	}
	if idx <= 0 {
		//预设下标超过元素集合范围，将下标设为最近元素的下标，此状态下为首元素下标
		idx = 0
	} else if idx >= len((*i.data))-1 {
		//预设下标超过元素集合范围，将下标设为最近元素的下标，此状态下为尾元素下标
		idx = len((*i.data)) - 1
	}
	if len((*i.data)) > 0 {
		//元素集合非空，迭代器下标设为预设下标
		i.index = idx
	} else {
		//元素集合为空，迭代器下标设为-1
		i.index = -1
	}
	//返回修改后的迭代器指针
	return i
}

func (i *Iterator) Value() (e interface{}) {
	if i == nil {
		//迭代器为nil，返回nil
		return nil
	}
	if len((*i.data)) == 0 {
		//元素集合为空，返回nil
		return nil
	}
	if i.index <= 0 {
		//下标超过元素集合范围下限，最近元素为首元素
		i.index = 0
	}
	if i.index >= len((*i.data)) {
		//下标超过元素集合范围上限，最近元素为尾元素
		i.index = len((*i.data)) - 1
	}
	//返回下标指向元素
	return (*i.data)[i.index]
}

func (i *Iterator) HasNext() (b bool) {
	if i == nil {
		//迭代器为nil时不能后移
		return false
	}
	if len((*i.data)) == 0 {
		//元素集合为空时不能后移
		return false
	}
	//下标到达元素集合上限时不能后移,否则可以后移
	return i.index < len((*i.data))
}

func (i *Iterator) Next() (b bool) {
	if i == nil {
		//迭代器为nil时返回false
		return false
	}
	if i.HasNext() {
		//满足后移条件时进行后移
		i.index++
		return true
	}
	if len((*i.data)) == 0 {
		//元素集合为空时下标设为-1同时返回false
		i.index = -1
		return false
	}
	//不满足后移条件时将下标设为尾元素下标并返回false
	i.index = len((*i.data)) - 1
	return false
}

func (i *Iterator) HasPre() (b bool) {
	if i == nil {
		//迭代器为nil时不能前移
		return false
	}
	if len((*i.data)) == 0 {
		//元素集合为空时不能前移
		return false
	}
	//下标到达元素集合范围下限时不能前移,否则可以后移
	return i.index >= 0
}

func (i *Iterator) Pre() (b bool) {
	if i == nil {
		//迭代器为nil时返回false
		return false
	}
	if i.HasPre() {
		//满足后移条件时进行前移
		i.index--
		return true
	}
	if len((*i.data)) == 0 {
		//元素集合为空时下标设为-1同时返回false
		i.index = -1
		return false
	}
	//不满足后移条件时将下标设为尾元素下标并返回false
	i.index = 0
	return false
}

func (i *Iterator) Display() {
	if i == nil {
		return
	}

	fmt.Println("Iterator Start...")
	for i := i.Begin(); i.HasNext(); i.Next() {
		fmt.Printf("idx:[%+v]\t value:[%+v]\n", i.index, i.Value())
	}
	fmt.Println()
}
