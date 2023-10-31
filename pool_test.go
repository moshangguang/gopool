package gopool

import (
	"math/rand"
	"testing"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

func testExample(t *testing.T, pool *WorkerPool) {
	intMap := make(map[int]struct{})
	forTimes := r.Intn(1000) + 1000 //循环次数最小为1000

	for i := 0; i < forTimes; i++ {
		intMap[rand.Int()] = struct{}{}
	}
	t.Logf("共有%d个数字待入队", len(intMap))
	ch := make(chan int, len(intMap))
	start := time.Now()
	for i := range intMap {
		num := i
		pool.Submit(func() {
			ch <- num
		})
	}
	t.Logf("入队耗时:%dms", time.Since(start).Milliseconds())
	start = time.Now()
	pool.StopWait()
	t.Logf("等待结束耗时:%dms", time.Since(start).Milliseconds())
	close(ch)
	workerCount := pool.WorkerCount()
	if workerCount == 0 {
		t.Log("协程池执行完毕后，协程数为0")
	} else {
		t.Fatalf("协程池执行完毕后，协程数不为0,workerCount:%d", workerCount)
	}
	for num := range ch {
		_, ok := intMap[num]
		if !ok {
			t.Fatalf("数字%d不在int map中", num)
		}
		delete(intMap, num)
	}
	if len(intMap) == 0 {
		t.Log("int map如预期清空")
	} else {
		t.Fatal("int map没有清空")
	}
}
func TestExample(t *testing.T) {
	pool := New(ArrayBlockQueueOption(128))
	testExample(t, pool)
	pool = New(LinkListBlockQueueOption(128))
	testExample(t, pool)
}

func TestWorkerPool_Stop(t *testing.T) {

}
func testMaxWorkerCountOption(t *testing.T, pool *WorkerPool, maxWorkerCount int) {
	if pool.maxWorkerCount != int32(maxWorkerCount) {
		t.Fatalf("协程池最大协程数与预期不符")
	} else {
		t.Log("协程池最大协程数与预期相等")
	}
	ch := make(chan struct{})
	for i := 0; i < maxWorkerCount; i++ {
		pool.Submit(func() {
			<-ch
		})
	}
	time.Sleep(time.Second)
	if pool.WorkerCount() != int32(maxWorkerCount) {
		t.Fatal("工作协程数没有达到最大协程数")
	} else {
		t.Logf("工作协程数达到最大协程数，workerCount：%d", pool.WorkerCount())
	}
	for i := 0; i < maxWorkerCount; i++ {
		ch <- struct{}{}
	}
	close(ch)
	pool.StopWait()
}
func TestMaxWorkerCountOption(t *testing.T) {
	maxWorkerCount := r.Intn(100) + 20
	pool := New(ArrayBlockQueueOption(128), MaxWorkerCountOption(maxWorkerCount))
	testMaxWorkerCountOption(t, pool, maxWorkerCount)
	maxWorkerCount = r.Intn(100) + 20
	pool = New(LinkListBlockQueueOption(128), MaxWorkerCountOption(maxWorkerCount))
	testMaxWorkerCountOption(t, pool, maxWorkerCount)
}
