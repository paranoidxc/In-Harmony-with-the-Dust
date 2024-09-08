package main

import (
	"container/heap"
	"fmt"
	"time"
)

type TimeSortedQueueItem struct {
	Time  int64
	Value interface{}
}

type TimeSortedQueue []*TimeSortedQueueItem

func (q TimeSortedQueue) Len() int           { return len(q) }
func (q TimeSortedQueue) Less(i, j int) bool { return q[i].Time < q[j].Time }
func (q TimeSortedQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *TimeSortedQueue) Push(v interface{}) {
	*q = append(*q, v.(*TimeSortedQueueItem))
}

func (q *TimeSortedQueue) Pop() interface{} {
	n := len(*q)
	item := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return item
}

func NewTimeSortedQueue(items ...*TimeSortedQueueItem) *TimeSortedQueue {
	q := make(TimeSortedQueue, len(items))
	for i, item := range items {
		q[i] = item
	}
	heap.Init(&q)
	return &q
}

func (q *TimeSortedQueue) PushItem(time int64, value interface{}) {
	heap.Push(q, &TimeSortedQueueItem{
		Time:  time,
		Value: value,
	})
}

func (q *TimeSortedQueue) PopItem() interface{} {
	if q.Len() == 0 {
		return nil
	}

	//return heap.Pop(q).(*TimeSortedQueueItem).Value
	return heap.Pop(q).(*TimeSortedQueueItem)
}

func main() {
	// 创建一个新的时间排序队列
	queue := NewTimeSortedQueue()

	// 向队列中添加一些项目
	queue.PushItem(time.Now().Unix(), "Event at the current time")
	queue.PushItem(time.Now().Add(2*time.Hour).Unix(), "Event in 2 hours")
	queue.PushItem(time.Now().Add(-1*time.Hour).Unix(), "Event 1 hour ago")

	// 按时间顺序弹出项目
	for queue.Len() > 0 {
		item := queue.PopItem()
		fmt.Printf("Popped item: %v at %v\n", item, time.Unix(item.(*TimeSortedQueueItem).Time, 0))
	}
}
