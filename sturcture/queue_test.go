package structure

import (
	"runtime"
	"testing"
	"time"
)
import "fmt"

// 这是一个 Queue 使用的示例
// 生产者:消费者 = n: 1
func TestQueue1(t *testing.T ){
	queue := New(1)
	// 创建1000 个写入协程
	for i:=0; i < 1000; i++{
		go func(num int){
			for j:=0; j < 100;j++{
				queue.Push(fmt.Sprintf("[%d] :hello world %d", num, j))
				queue.SingleUP(false, 1)
			}
		}(i)
	}
	
	for true {
		count := 0
		select{
		case <-queue.Single():
			for v := queue.Pop(); v!=nil; {
				t.Log(v.(string))
				count ++
				v = queue.Pop()
			}

			queue.SingleDown()
			if count >= 1000 * 100{
				return
			}
		}
	}

	return
}

// 生产者：消费者 = 1 :n 
func TestQueue2(t *testing.T ){
	cpuNum := runtime.NumCPU()
	queue := New(int32(cpuNum))
	for i:=0; i < cpuNum; i++{
		go func(num int){
			for true {
				select{
				case <-queue.Single():
					for v := queue.Pop(); v!=nil; {
						t.Logf("[%d]: %s", num, v.(string))
						v = queue.Pop()
					}
					queue.SingleDown()
					t.Logf("[%d]: more and more !!!", num)
				}
			}
		}(i)
	}

	for i:=0; i < 1000 * 100; i++{
		queue.Push(fmt.Sprintf("hello world %d", i))
		queue.SingleUP(false, uint8(cpuNum))
	}

	time.Sleep(time.Second * 10)
	t.Logf("cpu num %d",  cpuNum)
	return
}


// 生产者：消费者 = m : n
func TestQueue3(t *testing.T ){
	cpuNum := runtime.NumCPU()
	queue := New(int32(cpuNum))
	for i:=0; i < cpuNum; i++{
		go func(num int){
			for true {
				select{
				case <-queue.Single():
					for v := queue.Pop(); v!=nil; {
						t.Logf("[%d]: %s", num, v.(string))
						v = queue.Pop()
						time.Sleep(time.Nanosecond )
					}
					queue.SingleDown()
					t.Logf("[%d]: more and more !!!", num)
				}
			}
		}(i)
	}

	// 创建1000 个写入协程
	for i:=0; i < 100; i++{
		go func(num int){
			for j:=0; j < 100;j++{
				queue.Push(fmt.Sprintf("[%d] :hello world %d", num, j))
				queue.SingleUP(false, uint8(cpuNum))
			}
		}(i)
	}


	time.Sleep(time.Second * 20)
	t.Logf("cpu num %d",  cpuNum)
	return
}
