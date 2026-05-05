// Package worker provides a configurable worker pool that processes tasks concurrently.
//
// Basic usage:
//
//	p := worker.NewPool(3)   // create a pool with 3 workers
//	p.Start()               // launch the worker goroutines
//	p.Submit(myTask)        // send tasks into the pool
//	p.Stop()                // signal no more tasks and wait for workers to finish
//
//	for r := range p.Results() {
//	    // handle each result
//	}
package worker

import (
	"sync"

	"github.com/hasanulkabir-md/go-worker-pool/task"
)

// Job wraps a Task with a numeric identifier so callers can correlate results.
type Job struct {
	ID   int
	Task task.Task
}

// Pool manages a fixed number of worker goroutines that pull jobs from a shared
// channel and send their results to a results channel.
type Pool struct {
	numWorkers int
	jobs       chan Job
	results    chan task.Result
	wg         sync.WaitGroup
}

// NewPool returns a new Pool configured with the given number of workers.
// numWorkers must be greater than zero.
func NewPool(numWorkers int) *Pool {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	return &Pool{
		numWorkers: numWorkers,
		// Buffer the channels so that submitters and collectors do not have to
		// run in lock-step with the workers.
		jobs:    make(chan Job, numWorkers*2),
		results: make(chan task.Result, numWorkers*2),
	}
}

// Start launches the worker goroutines.  It must be called before any jobs are
// submitted.  Calling Start more than once on the same Pool is not supported.
func (p *Pool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.runWorker()
	}
}

// runWorker is the goroutine body for a single worker.
// It reads jobs from the jobs channel until it is closed, then decrements the
// WaitGroup counter so Stop knows all workers have finished.
func (p *Pool) runWorker() {
	defer p.wg.Done()
	for job := range p.jobs {
		value, err := job.Task.Process()
		p.results <- task.Result{
			TaskID: job.ID,
			Value:  value,
			Err:    err,
		}
	}
}

// Submit enqueues a job for processing.  It blocks if the internal jobs channel
// buffer is full.  Submit must not be called after Stop.
func (p *Pool) Submit(id int, t task.Task) {
	p.jobs <- Job{ID: id, Task: t}
}

// Stop signals that no more jobs will be submitted, waits for all in-flight
// jobs to complete, and then closes the results channel so that callers ranging
// over Results() will exit cleanly.
func (p *Pool) Stop() {
	close(p.jobs)
	// Wait in a separate goroutine so that the caller can still drain Results()
	// concurrently without deadlocking on a full results buffer.
	go func() {
		p.wg.Wait()
		close(p.results)
	}()
}

// Results returns the read-only channel on which completed results are published.
// Callers should range over this channel until it is closed (i.e. after Stop has
// been called and all workers have finished).
func (p *Pool) Results() <-chan task.Result {
	return p.results
}
