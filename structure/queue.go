package structure

// 参考网址：https://blog.csdn.net/qq_26981997/article/details/78773487

import (
	"sync/atomic"
	"unsafe"
)

// 单向链表
type Node struct {
	next *Node
	data interface{}
}

type Queue struct {
	head, tail    *Node
	signalCap     int32
	signalCounter int32
	notify        chan struct{}
}

func NewQueue(signalCap int32) *Queue {
	q := &Queue{}
	q.head = &Node{}
	q.tail = q.head

	q.signalCap = signalCap
	q.signalCounter = 0

	q.notify = make(chan struct{}, 2*signalCap)

	return q
}

func (qu *Queue) Single() <-chan struct{} {
	return qu.notify
}

func (qu *Queue) Push(x interface{}) {
	n := &Node{data: x}
	prev := (*Node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&qu.head)), unsafe.Pointer(n)))
	prev.next = n
}

func (qu *Queue) SingleUP(force bool) {
	if force {
		num := cap(qu.notify) - len(qu.notify) + 1
		for i := 0; i < num; i++ {
			qu.notify <- struct{}{}
		}
	} else if atomic.LoadInt32(&qu.signalCounter) < qu.signalCap {
		for atomic.LoadInt32(&qu.signalCounter) < qu.signalCap {
			qu.notify <- struct{}{}
			atomic.AddInt32(&qu.signalCounter, 1)
		}
	}
}

// Note: 这个函数要在处理线程处理完后再去调用
// 否则会造成写数据阻塞
func (qu *Queue) SingleDown() {
	atomic.AddInt32(&qu.signalCounter, -1)
}

func (qu *Queue) Pop() (v interface{}) {
	for {
		// 取尾结点指针
		tail := qu.tail

		// 读取next节点
		next := tail.next
		if next == nil {
			break
		}

		// 替换尾指针
		swap := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&qu.tail)),
			unsafe.Pointer(tail), unsafe.Pointer(next))
		if swap {
			v = next.data
			break
		}
	}

	return v
}

func (qu *Queue) Close() {
	close(qu.notify)
}
