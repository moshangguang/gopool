package main

import (
	"sync"
	"testing"
	"time"
)

func TestArrayBlockDeque_PopFrontWithTimeout(t *testing.T) {
	fn := func() {}
	maxLen := 3
	queue := NewArrayQueue(maxLen)
	for i := 0; i < maxLen; i++ {
		queue.PushBack(fn)
	}
	for i := 0; i < maxLen; i++ {
		task := queue.PopFrontWithTimeout(time.Millisecond)
		if task == nil {
			t.Fatal("获取任务失败")
		} else {
			t.Log("获取任务成功")
		}
	}
	task := queue.PopFrontWithTimeout(time.Millisecond)
	if task == nil {
		t.Log("队列为空后，预期获取任务失败，但获取任务成功")
	} else {
		t.Fatal("预期获取任务失败，但获取任务成功")
	}
	start := time.Now()
	go func() {
		time.Sleep(300 * time.Millisecond)
		queue.PushBack(fn)
	}()
	task = queue.PopFrontWithTimeout(500 * time.Millisecond)
	millis := time.Since(start).Milliseconds()
	if task != nil && millis > 200 {
		t.Logf("阻塞队列如预期获取到任务，且获取时长为:%dms", millis)
	}
}
func TestArrayBlockDeque_PushBack(t *testing.T) {
	fn := func() {}
	maxLen := 3
	queue := NewArrayQueue(maxLen)
	for i := 0; i < maxLen; i++ {
		ok := queue.PushBack(fn)
		if ok {
			t.Log("任务插入队列成功")
		} else {
			t.Fatal("任务插入队列失败")
		}
	}
	ok := queue.PushBack(fn)
	if ok {
		t.Fatal("队列已满后，任务插入队列成功")
	} else {
		t.Log("队列已满，任务插入到队列失败")
	}
	maxLen = 100
	queue = NewArrayQueue(maxLen)
	wg := new(sync.WaitGroup)
	wg.Add(maxLen)
	for i := 0; i < maxLen; i++ {
		go func() {
			defer wg.Done()
			queue.PushBack(fn)
		}()
	}
	wg.Wait()
	size := queue.Size()
	if size != maxLen {
		t.Fatalf("队列预期长度:%d，实际队列长度:%d", maxLen, size)
	} else {
		t.Logf("队列预期长度和实际队列长度一致:%d", size)
	}
}

func TestArrayBlockDeque_Size(t *testing.T) {
	fn := func() {}
	maxLen := 3
	queue := NewArrayQueue(maxLen)
	for i := 0; i < maxLen-1; i++ {
		queue.PushBack(fn)
	}
	if queue.Size() == maxLen-1 {
		t.Logf("队列预计长度和实际长度一致:%d", queue.Size())
	} else {
		t.Fatalf("队列预计长度:%d，队列实际长度:%d", maxLen-1, queue.Size())
	}
	queue.PushBack(fn)
	if queue.Size() == maxLen {
		t.Logf("队列预计长度和实际长度一致:%d", queue.Size())
	} else {
		t.Logf("队列预计长度:%d，队列实际长度:%d", maxLen-1, queue.Size())
	}

}
func TestArrayBlockDeque_Close(t *testing.T) {
	fn := func() {}
	maxLen := 3
	queue := NewArrayQueue(maxLen)
	for i := 0; i < maxLen-1; i++ {
		queue.PushBack(fn)
	}
	queue.Close()
	ok := queue.IsClose()
	if ok {
		t.Log("队列已关闭")
	} else {
		t.Fatal("队列未如预期关闭")
	}
	ok = queue.PushBack(fn)
	if ok {
		t.Fatal("队列已关闭，预期插入任务失败，实际插入成功")
	} else {
		t.Log("队列已关闭，任务如预期插入失败")
	}
}
