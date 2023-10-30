package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	coreWorkerCount int32             //核心协程数
	maxWorkerCount  int32             //最大协程数
	workerCount     *int32            //当前协程数
	idleTimeout     time.Duration     //空闲时间
	blockQueue      BlockQueue        //阻塞队列
	panicHandler    func()            //panic处理器
	rejectHandler   func(task func()) //拒绝策略

	workerQueue chan func() //工作队列
	wg          *sync.WaitGroup
	stopOnce    sync.Once
	stopChan    chan struct{}
}

func New(options ...WorkerPoolOption) *WorkerPool {
	pool := &WorkerPool{
		coreWorkerCount: DefaultCoreWorkerCount,
		workerCount:     new(int32),
		workerQueue:     make(chan func()),
		wg:              new(sync.WaitGroup),
		stopChan:        make(chan struct{}, 1),
	}
	if len(options) != 0 {
		for _, option := range options {
			if option != nil {
				option(pool)
			}
		}
	}
	pool.formatter()
	go pool.dispatch()

	return pool
}
func (p *WorkerPool) formatter() {
	if p.blockQueue == nil {
		p.blockQueue = NewArrayQueue(DefaultQueueLen)
	}
	if p.coreWorkerCount < 0 {
		p.coreWorkerCount = 0
	}
	if p.maxWorkerCount <= 0 {
		p.maxWorkerCount = DefaultMaxWorkerCount
	}
	if p.idleTimeout <= 0 || p.idleTimeout > MaxIdleTimeout {
		p.idleTimeout = MaxIdleTimeout
	}
	if p.coreWorkerCount > p.maxWorkerCount {
		p.maxWorkerCount = p.coreWorkerCount
	}

	if p.panicHandler == nil {
		RecoverPanicHandlerOption()(p)
	}
	if p.rejectHandler == nil {
		CallerRunRejectHandlerOption()(p)
	}
}
func (p *WorkerPool) wrapperTask(task func()) func() {
	return func() {
		defer p.panicHandler()
		task()
	}
}
func (p *WorkerPool) wrapperTaskWithDone(task func()) (func(), chan struct{}) {
	ch := make(chan struct{})
	return func() {
		defer close(ch)
		task = p.wrapperTask(task)
		task()
	}, ch
}
func (p *WorkerPool) Submit(task func()) {
	if task == nil {
		return
	}
	task = p.wrapperTask(task)
	if ok := p.blockQueue.PushBack(task); ok {
		return
	}
	p.rejectHandler(task)
}
func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	task, done := p.wrapperTaskWithDone(task)
	if ok := p.blockQueue.PushBack(task); ok {
		<-done
		return
	}
	p.rejectHandler(task)
}
func (p *WorkerPool) dispatch() {
	var task func()
loop:
	for {
		task = p.blockQueue.PopFrontWithTimeout(p.idleTimeout)
		if task != nil {
			select {
			case p.workerQueue <- task: //尝试将任务放入工作队列
			default:
				// 没有一个协程可以执行任务，则看看工作协程数是否小于最大协程数，如果小于的话则开启一个新的协程
				if atomic.LoadInt32(p.workerCount) < p.maxWorkerCount {
					p.wg.Add(1)
					atomic.AddInt32(p.workerCount, 1)
					go p.worker(task)
				} else { //如果工作协程数>最大协程数，则将任务放进工作协程的通道
					p.workerQueue <- task
				}
			}
		}
		task = nil
		if p.blockQueue.Size() != 0 {
			continue
		}
		select {
		case <-p.stopChan:
			break loop
		default:
			continue
		}
	}
	timer := time.NewTimer(p.idleTimeout)
	for atomic.LoadInt32(p.workerCount) != 0 {
		select {
		case p.workerQueue <- nil:
		case <-timer.C:
			timer.Reset(p.idleTimeout)
		}
	}
	p.wg.Done()

}
func (p *WorkerPool) release() bool {
	for {
		if p.blockQueue.Size() == 0 && p.blockQueue.IsClose() {
			return true
		}
		workerCount := atomic.LoadInt32(p.workerCount)
		if workerCount <= p.coreWorkerCount {
			return false
		}
		swapped := atomic.CompareAndSwapInt32(p.workerCount, workerCount, workerCount-1)
		if swapped {
			return true
		}
	}
}

func (p *WorkerPool) worker(task func()) {
	defer p.wg.Done()
	defer atomic.AddInt32(p.workerCount, -1)
	for task != nil {
		task()
		task = nil
		select {
		case task = <-p.workerQueue:
		case <-time.After(p.idleTimeout):
			if p.release() {
				return
			}
			task = <-p.workerQueue
		}
	}

}
func (p *WorkerPool) stop(wait bool) {
	p.stopOnce.Do(func() {
		p.blockQueue.Close()
		p.stopChan <- struct{}{}
	})
	if wait {
		p.wg.Done()
	}
}
func (p *WorkerPool) Stop() {
	p.stop(false)
}
func (p *WorkerPool) StopWait() {
	p.stop(true)
}
