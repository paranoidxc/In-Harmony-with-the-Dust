package deque

const (
	CAP = 10
)

type node struct {
	data  [CAP]interface{} //用于承载元素的固定数组
	begin int16            //该结点在前方添加结点的下标
	end   int16            //该结点在后方添加结点的下标
	pre   *node            //该结点的前一个结点
	next  *node            //该节点的后一个结点
}

type noder interface {
	nextNode() (m *node)                   //返回下一个结点
	preNode() (m *node)                    //返回上一个结点
	value() (es []interface{})             //返回该结点所承载的所有元素
	pushFront(e interface{}) (first *node) //在该结点头部添加一个元素,并返回新首结点
	pushBack(e interface{}) (last *node)   //在该结点尾部添加一个元素,并返回新尾结点
	popFront() (first *node)               //弹出首元素并返回首结点
	popBack() (last *node)                 //弹出尾元素并返回尾结点
	front() (e interface{})                //返回首元素
	back() (e interface{})                 //返回尾元素
}

func createFirst() (n *node) {
	return &node{
		data:  [CAP]interface{}{},
		begin: CAP - 1,
		end:   CAP,
		pre:   nil,
		next:  nil,
	}
}

func createLast() (n *node) {
	return &node{
		data:  [CAP]interface{}{},
		begin: -1,
		end:   0,
		pre:   nil,
		next:  nil,
	}
}

func (n *node) nextNode() (m *node) {
	if n == nil {
		return nil
	}

	return n.next
}

func (n *node) preNode() (m *node) {
	if n == nil {
		return nil
	}

	return n.pre
}

// 返回当前结点所承载的所有元素
func (n *node) value() (vals []interface{}) {
	vals = make([]interface{}, 0, 0)
	if n == nil {
		return vals
	}

	if n.begin > n.end {
		return vals
	}

	vals = n.data[n.begin+1 : n.end]

	return vals
}

func (n *node) pushFront(e interface{}) (first *node) {
	if n == nil {
		return n
	}

	if n.begin >= 0 {
		//该结点仍有空间可用于承载元素
		n.data[n.begin] = e
		n.begin--
		return n
	}

	//该结点无空间承载,创建新的首结点用于存放
	m := createFirst()
	m.data[m.begin] = e
	m.next = n
	n.pre = m
	m.begin--

	return m
}

func (n *node) pushBack(e interface{}) (last *node) {
	if n == nil {
		return n
	}

	if n.end < int16(len(n.data)) {
		//该结点仍有空间可用于承载元素
		n.data[n.end] = e
		n.end++
		return n
	}

	//该结点无空间承载,创建新的尾结点用于存放
	m := createLast()
	m.data[m.end] = e
	m.pre = n
	n.next = m
	m.end++

	return m
}

func (n *node) front() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data[n.begin+1]
}

func (n *node) popFront() (first *node) {
	if n == nil {
		return nil
	}

	if n.begin < int16(len(n.data))-2 {
		n.begin++
		n.data[n.begin] = nil
		return n
	}
	if n.next != nil {
		//清除该结点下一节点的前结点指针
		n.next.pre = nil
	}

	return n.next
}

func (n *node) back() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data[n.end-1]
}

func (n *node) popBack() (last *node) {
	if n == nil {
		return nil
	}
	if n.end > 1 {
		//该结点仍有承载元素
		n.end--
		n.data[n.end] = nil
		return n
	}
	if n.pre != nil {
		//清除该结点上一节点的后结点指针
		n.pre.next = nil
	}

	return n.pre
}
