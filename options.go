package gopool

import "time"

type WorkerPoolOption func(*WorkerPool)

func RejectHandlerOption(handler func(task func())) WorkerPoolOption {
	return func(pool *WorkerPool) {
		pool.rejectHandler = handler
	}
}
func CallerRunRejectHandlerOption() WorkerPoolOption {
	return func(pool *WorkerPool) {
		pool.rejectHandler = func(task func()) {
			task()
		}
	}
}

func RecoverPanicHandlerOption() WorkerPoolOption {
	return func(pool *WorkerPool) {
		pool.panicHandler = func() {
			v := recover()
			if v == nil {
				return
			}
			_ = v
		}
	}
}

func ArrayBlockQueueOption(maxLen int) WorkerPoolOption {
	if maxLen <= 0 {
		maxLen = DefaultQueueLen
	}
	return func(pool *WorkerPool) {
		pool.blockQueue = NewArrayQueue(maxLen)
	}
}

func LinkListBlockQueueOption(maxLen int) WorkerPoolOption {
	if maxLen <= 0 {
		maxLen = DefaultQueueLen
	}
	return func(pool *WorkerPool) {
		pool.blockQueue = NewLinkListBlockQueue(maxLen)
	}
}
func CoreWorkerCountOption(coreWorkerCount int) WorkerPoolOption {
	if coreWorkerCount < 0 {
		coreWorkerCount = DefaultCoreWorkerCount
	}
	return func(pool *WorkerPool) {
		pool.coreWorkerCount = int32(coreWorkerCount)
		if pool.maxWorkerCount < pool.coreWorkerCount {
			pool.maxWorkerCount = pool.coreWorkerCount
		}
	}
}

func MaxWorkerCountOption(maxWorkerCount int) WorkerPoolOption {
	if maxWorkerCount <= 0 {
		maxWorkerCount = DefaultMaxWorkerCount
	}
	return func(pool *WorkerPool) {
		if pool.coreWorkerCount < int32(maxWorkerCount) {
			pool.maxWorkerCount = int32(maxWorkerCount)
		}
	}
}

func IdleTimeoutOption(idleTimeout time.Duration) WorkerPoolOption {
	if idleTimeout <= 0 {
		idleTimeout = MaxIdleTimeout
	}
	return func(pool *WorkerPool) {
		pool.idleTimeout = idleTimeout
	}
}
