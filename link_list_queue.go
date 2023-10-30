package gopool

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

type LinklistBlockQueue struct {
	BaseQueue
	list   *LinkList
	mu     *sync.RWMutex
	ch     chan func()
	weight *semaphore.Weighted
}

func (queue *LinklistBlockQueue) popFront(ctx context.Context) (task func()) {
	err := queue.weight.Acquire(ctx, 1)
	if err != nil {
		return
	}
	queue.mu.Lock()
	defer queue.mu.Unlock()
	task = queue.list.PopFront()
	return
}
func (queue *LinklistBlockQueue) PopFrontWithTimeout(timeout time.Duration) (task func()) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	task = queue.popFront(ctx)
	return
}

func (queue *LinklistBlockQueue) Size() int {
	queue.mu.RLock()
	defer queue.mu.RUnlock()
	return queue.list.Size()
}

func NewLinkListBlockQueue(maxLen int) BlockQueue {
	weighted := semaphore.NewWeighted(int64(maxLen))
	ok := weighted.TryAcquire(int64(maxLen))
	if !ok {
		panic("error NewWeighted")
	}
	return &LinklistBlockQueue{
		BaseQueue: NewBaseQueue(),
		list:      NewLinkList(maxLen),
		mu:        new(sync.RWMutex),
		weight:    weighted,
	}
}
func (queue *LinklistBlockQueue) PopFront() (task func()) {
	task = queue.popFront(context.Background())
	return
}
func (queue *LinklistBlockQueue) PushBack(task func()) bool {
	if task == nil {
		return false
	}
	if queue.IsClose() {
		return false
	}
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if ok := queue.list.PushBack(task); !ok {
		return false
	}
	queue.weight.Release(1)
	return true
}

type Element struct {
	task func()
	next *Element
}

type LinkList struct {
	maxLen int
	len    int
	head   *Element
}

func NewLinkList(maxLen int) *LinkList {
	return &LinkList{
		len:    0,
		maxLen: maxLen,
		head:   nil,
	}
}

func (s *LinkList) PushBack(task func()) bool {
	if s.len >= s.maxLen {
		return false
	}
	elem := &Element{
		task: task,
		next: nil,
	}
	if s.head == nil {
		s.head = elem
	} else {
		current := s.head
		for current.next != nil {
			current = current.next
		}
		current.next = elem
	}
	s.len++
	return true
}

func (s *LinkList) PopFront() (task func()) {
	if s.head == nil {
		return
	}
	task = s.head.task
	s.head = s.head.next
	s.len--
	return
}

func (s *LinkList) Size() int {
	return s.len
}
