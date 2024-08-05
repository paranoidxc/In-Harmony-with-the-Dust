package link

type node struct {
	data interface{} //结点所承载的元素
	pre  *node       //前结点指针
	next *node       //后结点指针
}

type noder interface {
	preNode() (m *node)     //返回前结点指针
	nextNode() (m *node)    //返回后结点指针
	insertPre(pre *node)    //在该结点前插入结点并建立连接
	insertNext(next *node)  //在该结点后插入结点并建立连接
	erase()                 //删除该结点,并使该结点前后两结点建立连接
	value() (e interface{}) //返回该结点所承载的元素
	setValue(e interface{}) //修改该结点承载元素为e
}

func newNode(e interface{}) (n *node) {
	return &node{
		data: e,
		pre:  nil,
		next: nil,
	}
}

func (n *node) preNode() (pre *node) {
	if n == nil {
		return
	}
	return n.pre
}

func (n *node) nextNode() (next *node) {
	if n == nil {
		return
	}
	return n.next
}

func (n *node) insertPre(pre *node) {
	if n == nil || pre == nil {
		return
	}
	pre.next = n
	pre.pre = n.pre
	if n.pre != nil {
		n.pre.next = pre
	}
	n.pre = pre
}

func (n *node) insertNext(next *node) {
	if n == nil || next == nil {
		return
	}
	next.pre = n
	next.next = n.next
	if n.next != nil {
		n.next.pre = next
	}
	n.next = next
}

func (n *node) erase() {
	if n == nil {
		return
	}
	if n.pre == nil && n.next == nil {
		return
	} else if n.pre == nil {
		n.next.pre = nil
	} else if n.next == nil {
		n.pre.next = nil
	} else {
		n.pre.next = n.next
		n.next.pre = n.pre
	}
	n = nil
}

func (n *node) value() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data
}

func (n *node) setValue(e interface{}) {
	if n == nil {
		return
	}
	n.data = e
}
