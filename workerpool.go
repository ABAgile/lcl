package lcl

import (
	"context"
	"log/slog"
	"sync"
)

type HandlerFunc[T any] func(ctx context.Context, job T, logger *slog.Logger)

// WorkerPool manages a pool of goroutines to process jobs concurrently.
type WorkerPool[T any] struct {
	numWorkers int
	jobs       <-chan T // The pool receives jobs from this channel.
	handler    HandlerFunc[T]
	logger     *slog.Logger
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool[T any](numWorkers int, jobs <-chan T, handler HandlerFunc[T], logger *slog.Logger) *WorkerPool[T] {
	return &WorkerPool[T]{
		numWorkers: numWorkers,
		jobs:       jobs,
		handler:    handler,
		logger:     logger,
	}
}

// Run starts the workers and blocks until jobs channel been closed
func (p *WorkerPool[T]) Run(ctx context.Context) {
	for i := range p.numWorkers {
		p.wg.Add(1)
		go func(ctx context.Context, id int) {
			defer p.wg.Done()
			logger := p.logger.With("worker_id", id)
			logger.Info("Worker started.")
			// internal goroutine function that processes jobs from the channel
			// it runs continuously until jobs channel is closed
			for job := range p.jobs {
				p.handler(ctx, job, logger)
			}
			logger.Info("Worker stopped.")

		}(ctx, i+1)
	}
	p.logger.Info("Worker pool started", "num_workers", p.numWorkers)
	p.wg.Wait()
	p.logger.Info("Worker pool stopped.")
}

func RunConcurrent[T any](jobs []T, workerCount int, fn func(context.Context, T)) {
	jobCh := make(chan T, workerCount*3) // buffered channel, improves throughput
	go func() {
		for _, j := range jobs {
			jobCh <- j
		}
		close(jobCh)
	}()
	handler := func(ctx context.Context, job T, logger *slog.Logger) {
		fn(ctx, job)
	}
	wp := NewWorkerPool(workerCount, jobCh, handler, slog.New(slog.DiscardHandler))
	wp.Run(context.Background())
}

func MapConcurrent[T any, R any](jobs []T, workerCount int, fn func(context.Context, T) R) []R {
	jobCh := make(chan Tuple2[int, T], workerCount*3) // buffered channel, improves throughput
	go func() {
		for i, j := range jobs {
			jobCh <- T2(i, j)
		}
		close(jobCh)
	}()
	results := make([]R, len(jobs))
	handler := func(ctx context.Context, job Tuple2[int, T], logger *slog.Logger) {
		i, j := job.Unpack()
		results[i] = fn(ctx, j)
	}
	wp := NewWorkerPool(workerCount, jobCh, handler, slog.New(slog.DiscardHandler))
	wp.Run(context.Background())
	return results
}
