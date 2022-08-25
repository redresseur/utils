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
	q := &Queue{head: &Node{}}
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
	(*Node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&qu.head)), unsafe.Pointer(n))).next = n
}

func (qu *Queue) SingleUP(force bool) {
	if !force {
		for atomic.LoadInt32(&qu.signalCounter) < qu.signalCap {
			qu.notify <- struct{}{}
			atomic.AddInt32(&qu.signalCounter, 1)
		}
	} else {
		num := cap(qu.notify) - len(qu.notify) + 1
		for i := 0; i < num; i++ {
			qu.notify <- struct{}{}
		}
	}
}

// Note: 这个函数要在处理线程处理完后再去调用
// 否则会造成写数据阻塞
func (qu *Queue) SingleDown() {
	atomic.AddInt32(&qu.signalCounter, -1)
}

func (qu *Queue) Pop() interface{} {
	for {
		if tail, next := qu.tail, qu.tail.next; next == nil { // 取尾结点指针, 读取next节点
			break
		} else if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&qu.tail)),
			unsafe.Pointer(tail), unsafe.Pointer(next)) { // 替换尾指针
			return next.data
		}
	}
	return nil
}

func (qu *Queue) Close() {
	close(qu.notify)
}
