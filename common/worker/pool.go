package worker

import (
	"sync"
)

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func NewPool(maxQueueSize int) *Pool {
	p := Pool{
		work: make(chan Worker,maxQueueSize),
	}
	return &p
}

func (p *Pool) Add(w Worker) {
	//if _, isClose := <-p.work;!isClose{
	//	return
	//}
	p.wg.Add(1)
	p.work <- w
}

func (p *Pool) Run(MaxPushWorker int) {
	for i := 0; i < MaxPushWorker; i++ {
		go func() {
			for w := range p.work {
				w.Run()
				p.wg.Done()
			}
		}()
	}
}

func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}

