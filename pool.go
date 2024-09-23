package main

import (
	"sync"
)

type Job func()

type Pool struct {
	workQueue chan Job
	wg        sync.WaitGroup
}

func (p *Pool) AddJob(job Job) {
	p.workQueue <- job
}

func (p *Pool) Wait() {
	close(p.workQueue)
	p.wg.Wait()
}

func NewPool(workerCount int) *Pool {
	pool := &Pool{
		workQueue: make(chan Job),
	}

	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer pool.wg.Done()
			for job := range pool.workQueue {
				job()
			}
		}()
	}
	return pool
}
