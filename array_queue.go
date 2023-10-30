package gopool

import "time"

type ArrayBlockDeque struct {
	BaseQueue
	ch chan func()
}

func (queue *ArrayBlockDeque) PopFrontWithTimeout(timeout time.Duration) (task func()) {
	select {
	case task = <-queue.ch:
		return
	case <-time.After(timeout):
		return
	}
}

func (queue *ArrayBlockDeque) PopFront() (task func()) {
	task = <-queue.ch
	return
}

func (queue *ArrayBlockDeque) PushBack(task func()) (ok bool) {
	if queue.IsClose() {
		return
	}
	select {
	case queue.ch <- task:
		return true
	default:
		return false
	}
}

func (queue *ArrayBlockDeque) Size() int {
	return len(queue.ch)
}

func NewArrayQueue(maxLen int) BlockQueue {
	return &ArrayBlockDeque{
		BaseQueue: NewBaseQueue(),
		ch:        make(chan func(), maxLen),
	}
}
