package worker

import (
	"sync"
)

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func NewPool() *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	return &p
}

func (p *Pool) Add(w Worker) {
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
