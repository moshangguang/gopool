package main

import "time"

const (
	DefaultCoreWorkerCount = 1 << 4
	DefaultMaxWorkerCount  = 1 << 10
	MaxIdleTimeout         = 5 * time.Second
	DefaultQueueLen        = 1 << 7
	QueueClose             = 1
)
