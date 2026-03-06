package workerpool

import "sync"

type Task func()

type Pool struct {
	workers int
	tasks   chan Task
	wg      sync.WaitGroup
}

func New(workers int) *Pool {
	p := &Pool{
		workers: workers,
		tasks:   make(chan Task),
	}

	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	for task := range p.tasks {
		task()
		p.wg.Done()
	}
}

func (p *Pool) Submit(task Task) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

// Usage in Kafka consumer

// pool := workerpool.New(10)

// for {
//     msg, _ := reader.ReadMessage(ctx)

//     pool.Submit(func() {
//         process(msg)
//     })
// }
