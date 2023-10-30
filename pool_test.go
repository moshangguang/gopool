package gopool

import "testing"

func TestWorkerPool_StopWait(t *testing.T) {
	texts := []string{"a", "b", "c", "d", "e"}
	ch := make(chan string, len(texts))
	pool := New()
	for _, text := range texts {
		item := text
		pool.Submit(func() {
			ch <- item
		})
	}
	pool.StopWait()
	close(ch)
	workerCount := pool.WorkerCount()
	if workerCount == 0 {
		t.Log("协程池执行完毕后，协程数为0")
	} else {
		t.Fatalf("协程池执行完毕后，协程数不为0,workerCount:%d", workerCount)
	}
	textMap := map[string]struct{}{}
	for item := range ch {
		textMap[item] = struct{}{}
	}
	if len(textMap) < len(texts) {
		t.Fatalf("从通道获取到的字符串数量缺失,len(textMap):%d", len(textMap))
	} else {
		t.Log("从通道获取到符合预期数量的字符串")
	}
	for _, req := range texts {
		if _, ok := textMap[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}
