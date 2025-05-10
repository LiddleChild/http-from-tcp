package worker

import (
	"sync"
)

type Pool[T any] struct {
	wg      sync.WaitGroup
	ch      chan T
	handler func(T)
}

func NewPool[T any](size int) *Pool[T] {
	p := &Pool[T]{
		ch:      make(chan T, size),
		handler: nil,
	}

	for range size {
		go p.worker()
	}

	return p
}

func (p *Pool[T]) TaskFunc(handler func(T)) {
	p.handler = handler
}

func (p *Pool[T]) Close() {
	close(p.ch)
	p.wg.Wait()
}

func (p *Pool[T]) QueueTask(task T) {
	p.ch <- task
}

func (p *Pool[T]) worker() {
	p.wg.Add(1)
	for {
		task, ok := <-p.ch
		if !ok {
			break
		}

		p.handler(task)
	}
	p.wg.Done()
}
