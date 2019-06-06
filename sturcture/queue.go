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
	head , tail *Node
	signalCap int32
	signalCounter int32
	notify chan struct{}
}

func NewQueue(signalCap int32) *Queue {
	q := &Queue{}
	q.head = &Node{}
	q.tail = q.head

	q.signalCap = signalCap
	q.signalCounter = 0

	q.notify = make(chan struct{}, 2* signalCap)

	return q
}

func (qu *Queue)Single()<-chan struct{}{
	return qu.notify
}

func (qu *Queue)Push(x interface{}){
	n := &Node{data: x}
	prev := (*Node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&qu.head)), unsafe.Pointer(n)))
	prev.next = n
}

func (qu *Queue)SingleUP(force bool){
	if force{
		qu.notify <- struct{}{}
	} else if atomic.LoadInt32(&qu.signalCounter) < qu.signalCap{
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

func (qu *Queue)Pop()interface{}  {
	tail := qu.tail
	next := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))) // acquire
	if next != nil {
		qu.tail = next
		v := next.data
		next.data = nil
		return v
	}

	return nil
}

func (qu *Queue)Close()  {
	close(qu.notify)
}