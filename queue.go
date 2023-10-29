package main

import "time"

type BlockQueue interface {
	PopFrontWithTimeout(timeout time.Duration) (task func())
	PushBack(task func()) (ok bool)
	Size() int
	Close()
	IsClose() bool
}
