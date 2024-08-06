package priority_queue

import (
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
	"sync"
)

type PriorityQueue struct {
	data  []interface{}         //动态数组
	len   uint64                //实际使用长度
	cap   uint64                //实际占用的空间的容量
	cmp   comparator.Comparator //该优先队列的比较器
	mutex sync.Mutex            //并发控制锁
}

type priorityQueueer interface {
	Size() (num uint64)   //返回该容器存储的元素数量
	Clear()               //清空该容器
	Empty() (b bool)      //判断该容器是否为空
	Push(e interface{})   //将元素e插入该容器
	Pop()                 //弹出顶部元素
	Top() (e interface{}) //返回顶部元素
}

func New(cmps ...comparator.Comparator) (pq *PriorityQueue) {
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	//比较器为nil时后续的增删将会去寻找默认比较器
	return &PriorityQueue{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

func (pq *PriorityQueue) Size() (num uint64) {
	if pq == nil {
		pq = New()
	}
	return pq.len
}

func (pq *PriorityQueue) Clear() {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	//清空已分配的空间
	pq.data = make([]interface{}, 1, 1)
	pq.len = 0
	pq.cap = 1
}

func (pq *PriorityQueue) Empty() bool {
	if pq == nil {
		pq = New()
	}
	return pq.len == 0
}

func (pq *PriorityQueue) Push(e interface{}) {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	//判断是否存在比较器,不存在则寻找默认比较器,若仍不存在则直接结束
	if pq.cmp == nil {
		pq.cmp = comparator.GetCmp(e)
	}
	if pq.cmp == nil {
		return
	}
	//先判断是否需要扩容,同时使用和vector相同的扩容策略
	//即先翻倍扩容再固定扩容,随后在末尾插入元素e
	if pq.len < pq.cap {
		//还有冗余,直接添加
		pq.data[pq.len] = e
	} else {
		//冗余不足,需要扩容
		if pq.cap <= 65536 {
			//容量翻倍
			if pq.cap == 0 {
				pq.cap = 1
			}
			pq.cap *= 2
		} else {
			//容量增加2^16
			pq.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
		pq.data[pq.len] = e
	}
	pq.len++
	//到此时,元素以插入到末尾处,同时插入位的元素的下标为pq.len-1,随后将对该位置的元素进行上升
	//即通过比较它逻辑上的父结点进行上升
	pq.up(pq.len - 1)
}

func (pq *PriorityQueue) up(p uint64) {
	if p == 0 {
		//以及上升到顶部,直接结束即可
		return
	}
	if pq.cmp(pq.data[(p-1)/2], pq.data[p]) > 0 {
		//判断该结点和其父结点的关系
		//满足给定的比较函数的关系则先交换该结点和父结点的数值,随后继续上升即可
		pq.data[p], pq.data[(p-1)/2] = pq.data[(p-1)/2], pq.data[p]
		pq.up((p - 1) / 2)
	}
}

func (pq *PriorityQueue) Pop() {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return
	}
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	//将最后一位移到首位,随后删除最后一位,即删除了首位,同时判断是否需要缩容
	pq.data[0] = pq.data[pq.len-1]
	pq.len--
	//缩容判断,缩容策略同vector,即先固定缩容在折半缩容
	if pq.cap-pq.len >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		pq.cap -= 65536
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	} else if pq.len*2 < pq.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		pq.cap /= 2
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	}
	//判断是否为空,为空则直接结束
	if pq.Empty() {
		return
	}
	//对首位进行下降操作,即对比其逻辑上的左右结点判断是否应该下降,再递归该过程即可
	pq.down(0)
}

func (pq *PriorityQueue) down(p uint64) {
	q := p
	//先判断其左结点是否在范围内,然后在判断左结点是否满足下降条件
	if 2*p+1 <= pq.len-1 && pq.cmp(pq.data[p], pq.data[2*p+1]) > 0 {
		q = 2*p + 1
	}
	//在判断右结点是否在范围内,同时若判断右节点是否满足下降条件
	if 2*p+2 <= pq.len-1 && pq.cmp(pq.data[q], pq.data[2*p+2]) > 0 {
		q = 2*p + 2
	}
	//根据上面两次判断,从最小一侧进行下降
	if p != q {
		//进行交互,递归下降
		pq.data[p], pq.data[q] = pq.data[q], pq.data[p]
		pq.down(q)
	}
}

func (pq *PriorityQueue) Top() (e interface{}) {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return nil
	}
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	e = pq.data[0]
	return e
}
