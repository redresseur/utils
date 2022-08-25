package structure

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

// 这是一个 Queue 使用的示例
// 生产者:消费者 = n: 1
func TestQueue1(t *testing.T) {
	queue := NewQueue(1)
	group := sync.WaitGroup{}
	// 创建1000 个写入协程
	group.Add(1000 * 100)
	for i := 0; i < 1000; i++ {
		go func(num int) {
			for j := 0; j < 100; j++ {
				queue.Push(fmt.Sprintf("[%d] :hello world %d", num, j))
				queue.SingleUP(false)
			}
		}(i)
	}

	go func() {
		for {
			select {
			case _, ok := <-queue.Single():
				if !ok {
					return
				}
			}

			for {
				if v := queue.Pop(); v != nil {
					fmt.Println(v.(string))
					group.Done()
				} else {
					break
				}
			}

			queue.SingleDown()
		}
	}()

	group.Wait()
	queue.SingleUP(true)
	queue.Close()

	return
}

// 生产者：消费者 = 1 :n
func TestQueue2(t *testing.T) {
	cpuNum := runtime.NumCPU()
	queue := NewQueue(int32(cpuNum))
	group := sync.WaitGroup{}

	for i := 0; i < cpuNum; i++ {
		go func(num int) {
			for {
				select {
				case _, ok := <-queue.Single():
					if !ok {
						return
					}
				}

				for {
					if v := queue.Pop(); v != nil {
						fmt.Printf("[%d]: %s\n", num, v.(string))
						group.Done()
					} else {
						break
					}
				}

				queue.SingleDown()
				fmt.Printf("[%d]: more and more !!!\n", num)
			}
		}(i)
	}

	group.Add(1000 * 100)
	for i := 0; i < 1000*100; i++ {
		queue.Push(fmt.Sprintf("hello world %d", i))
		queue.SingleUP(false)
	}

	group.Wait()
	queue.Close()
	t.Logf("cpu num %d", cpuNum)
	return
}

// 生产者：消费者 = m : n
func TestQueue3(t *testing.T) {
	cpuNum := runtime.NumCPU()
	queue := NewQueue(int32(cpuNum))
	group := sync.WaitGroup{}

	for i := 0; i < cpuNum; i++ {
		go func(num int) {
			for {
				select {
				case _, ok := <-queue.Single():
					if !ok {
						return
					}
				}

				for v := queue.Pop(); v != nil; v = queue.Pop() {
					fmt.Printf("[%d]: %s\n", num, v.(string))
					group.Done()
				}

				queue.SingleDown()
				fmt.Printf("[%d]: more and more !!!\n", num)
			}

		}(i)
	}

	group.Add(1000 * 100)
	// 创建1000 个写入协程
	for i := 0; i < 1000; i++ {
		go func(num int) {
			for j := 0; j < 100; j++ {
				queue.Push(fmt.Sprintf("[%d] :hello world %d", num, j))
				queue.SingleUP(false)
			}
		}(i)
	}

	group.Wait()
	queue.Close()
	t.Logf("cpu num %d", cpuNum)
	return
}

func BenchmarkQueue_Push(b *testing.B) {
	var queue *Queue
	once := sync.Once{}
	once.Do(func() {
		queue = NewQueue(1)
		go func() {
			for {
				select {
				case _, ok := <-queue.Single():
					if !ok {
						return
					}
				}

				for {
					if v := queue.Pop(); v != nil {
					} else {
						break
					}
				}

				queue.SingleDown()
			}
		}()
	})

	for i := 0; i < b.N; i++ {
		queue.Push(i)
		queue.SingleUP(false)
	}
}

// 只使用channel
func BenchmarkQueue_Push1(b *testing.B) {
	var queue chan interface{}
	once := sync.Once{}
	once.Do(func() {
		queue = make(chan interface{}, 1)
		go func() {
			for {
				select {
				case _, ok := <-queue:
					if !ok {
						return
					}
				}
			}
		}()
	})

	for i := 0; i < b.N; i++ {
		queue <- i
	}
}

func BenchmarkQueue_Push2(b *testing.B) {
	var queue *Queue
	once := sync.Once{}
	once.Do(func() {
		queue = NewQueue(4)
		for i := 0; i < 4; i++ {
			go func() {
				for {
					select {
					case _, ok := <-queue.Single():
						if !ok {
							return
						}
					}

					for {
						if v := queue.Pop(); v != nil {
						} else {
							break
						}
					}

					queue.SingleDown()
				}
			}()
		}
	})

	for i := 0; i < b.N; i++ {
		queue.Push(i)
		queue.SingleUP(false)
	}
}

// 只使用channel
func BenchmarkQueue_Push3(b *testing.B) {
	var queue chan interface{}
	once := sync.Once{}
	once.Do(func() {
		queue = make(chan interface{}, 4)
		for i := 0; i < 4; i++ {
			go func() {
				for {
					select {
					case _, ok := <-queue:
						if !ok {
							return
						}
					}
				}
			}()
		}
	})

	for i := 0; i < b.N; i++ {
		queue <- i
	}
}

func BenchmarkPushAndPop(b *testing.B) {
	queue := NewQueue(1)
	for i := 0; i < b.N; i++ {
		queue.Push(1)
		queue.Pop()
	}
	queue.Close()
}

func BenchmarkPushAndPopConcurrency(b *testing.B) {
	queue := NewQueue(1)
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			queue.Push(1)
			queue.Pop()
		}
	})
	queue.Close()
}
