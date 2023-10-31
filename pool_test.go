package gopool

import (
	"math/rand"
	"testing"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

func TestExample(t *testing.T) {
	intMap := make(map[int]struct{})
	forTimes := r.Intn(1000) + 1000 //循环次数最小为1000

	for i := 0; i < forTimes; i++ {
		intMap[rand.Int()] = struct{}{}
	}
	t.Logf("共有%d个数字待入队", len(intMap))
	ch := make(chan int, len(intMap))
	pool := New(MaxWorkerCountOption(2000))
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
