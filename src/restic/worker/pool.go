package worker

import (
	"fmt"
	"sync"
)

// Job is one unit of work. It is given to a Func, and the returned result and
// error are stored in Result and Error.
type Job struct {
	Data   interface{}
	Result interface{}
	Error  error
}

// Func does the actual work within a Pool.
type Func func(job Job, done <-chan struct{}) (result interface{}, err error)

// Pool implements a worker pool.
type Pool struct {
	f     Func
	done  chan struct{}
	wg    *sync.WaitGroup
	jobCh <-chan Job
	resCh chan<- Job
}

// New returns a new worker pool with n goroutines, each running the function
// f. The workers are started immediately.
func New(n int, f Func, jobChan <-chan Job, resultChan chan<- Job) *Pool {
	p := &Pool{
		f:     f,
		done:  make(chan struct{}),
		wg:    &sync.WaitGroup{},
		jobCh: jobChan,
		resCh: resultChan,
	}

	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.runWorker(i)
	}

	return p
}

// runWorker runs a worker function.
func (p *Pool) runWorker(numWorker int) {
	defer p.wg.Done()

	var (
		// enable the input channel when starting up a new goroutine
		inCh = p.jobCh
		// but do not enable the output channel until we have a result
		outCh chan<- Job

		job Job
		ok  bool
	)

	for {
		select {
		case <-p.done:
			return

		case job, ok = <-inCh:
			if !ok {
				fmt.Printf("in channel closed, worker exiting\n")
				return
			}

			job.Result, job.Error = p.f(job, p.done)
			inCh = nil
			outCh = p.resCh

		case outCh <- job:
			outCh = nil
			inCh = p.jobCh
		}
	}
}

// Cancel signals termination to all worker goroutines.
func (p *Pool) Cancel() {
	close(p.done)
}

// Wait waits for all worker goroutines to terminate, afterwards the output
// channel is closed.
func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.resCh)
}
