package reactor

import (
	"sync"
	"runtime"
	"fmt"
)

type OnDone func ()

type Task interface {
	Do()
}

type Pool interface {
	Start()
	Stop()
	PostTask(task Task)
}

type worker struct {
	taskCh chan Task
	quitCh chan *sync.WaitGroup
}

type pool struct {
	taskCh chan Task
	workers []*worker

	lock sync.Mutex
	running bool
}

func newWorker(taskCh chan Task) *worker {
	return &worker{
		taskCh: taskCh,
		quitCh: make(chan *sync.WaitGroup),
	}
}

func (this *worker) start() {
	go func() {
		for {
			select {
			case task := <- this.taskCh:
				task.Do()
			case wg := <- this.quitCh:
				wg.Done()
				break
			}
		}
	}()
}

func (this *worker) stop(wg *sync.WaitGroup) {
	this.quitCh <- wg
}

func NewPool(nWorkers int) Pool {
	if nWorkers == 0 {
		nWorkers = runtime.NumCPU() / 2
		if nWorkers == 0 {
			nWorkers = 1
		}
	}
	fmt.Println("nWorkers:", nWorkers)
	ret := &pool{
		taskCh: make(chan Task),
		workers: make([]*worker, nWorkers),
	}

	for i := 0; i < nWorkers; i++ {
		ret.workers[i] = newWorker(ret.taskCh)
	}

	return ret
}

func (this *pool) Start() {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.running {
		return
	}
	this.running = true

	for _, worker := range this.workers {
		worker.start()
	}
}

func (this *pool) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !this.running {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(this.workers))
	for _, worker := range this.workers {
		worker.stop(&wg)
	}
	wg.Wait()
	this.running = false
}

func (this *pool) PostTask(task Task) {
	this.taskCh <- task
}
