package main

import "sync/atomic"

type BaseQueue struct {
	close *int32
}

func (queue *BaseQueue) IsClose() bool {
	return atomic.LoadInt32(queue.close) == QueueClose
}
func (queue *BaseQueue) Close() {
	atomic.StoreInt32(queue.close, QueueClose)
}
func NewBaseQueue() BaseQueue {
	return BaseQueue{
		close: new(int32),
	}
}
