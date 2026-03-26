package lcl

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- WorkerPool ---

func TestWorkerPool_Run(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		jobs       []int
		wantSum    int
	}{
		{"single worker", 1, []int{1, 2, 3, 4, 5}, 15},
		{"multiple workers", 4, []int{1, 2, 3, 4, 5}, 15},
		{"more workers than jobs", 10, []int{1, 2, 3}, 6},
		{"empty jobs", 3, []int{}, 0},
		{"single job", 2, []int{42}, 42},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jobCh := make(chan int, len(tc.jobs)+1)
			for _, j := range tc.jobs {
				jobCh <- j
			}
			close(jobCh)

			var mu sync.Mutex
			sum := 0
			handler := func(_ context.Context, job int, _ *slog.Logger) {
				mu.Lock()
				sum += job
				mu.Unlock()
			}

			wp := NewWorkerPool(tc.numWorkers, jobCh, handler, slog.New(slog.DiscardHandler))
			wp.Run(context.Background())

			assert.Equal(t, tc.wantSum, sum)
		})
	}
}

func TestWorkerPool_Run_ContextPropagated(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "sentinel")

	jobCh := make(chan int, 1)
	jobCh <- 1
	close(jobCh)

	var got any
	handler := func(ctx context.Context, _ int, _ *slog.Logger) {
		got = ctx.Value(ctxKey{})
	}

	wp := NewWorkerPool(1, jobCh, handler, slog.New(slog.DiscardHandler))
	wp.Run(ctx)

	assert.Equal(t, "sentinel", got)
}

// --- RunConcurrent ---

func TestRunConcurrent(t *testing.T) {
	tests := []struct {
		name        string
		jobs        []int
		workerCount int
		wantCount   int32
	}{
		{"all jobs processed", []int{1, 2, 3, 4, 5}, 2, 5},
		{"single worker", []int{10, 20}, 1, 2},
		{"empty input", []int{}, 3, 0},
		{"more workers than jobs", []int{1}, 10, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var count atomic.Int32
			RunConcurrent(tc.jobs, tc.workerCount, func(_ context.Context, _ int) {
				count.Add(1)
			})
			assert.Equal(t, tc.wantCount, count.Load())
		})
	}
}

func TestRunConcurrent_AllJobsRun(t *testing.T) {
	jobs := []string{"a", "b", "c", "d", "e"}
	var mu sync.Mutex
	var seen []string

	RunConcurrent(jobs, 3, func(_ context.Context, s string) {
		mu.Lock()
		seen = append(seen, s)
		mu.Unlock()
	})

	sort.Strings(seen)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, seen)
}

// --- MapConcurrent ---

func TestMapConcurrent(t *testing.T) {
	tests := []struct {
		name        string
		jobs        []int
		workerCount int
		mapper      func(context.Context, int) int
		want        []int
	}{
		{
			name:        "doubles each element",
			jobs:        []int{1, 2, 3, 4, 5},
			workerCount: 3,
			mapper:      func(_ context.Context, n int) int { return n * 2 },
			want:        []int{2, 4, 6, 8, 10},
		},
		{
			name:        "identity",
			jobs:        []int{7, 8, 9},
			workerCount: 2,
			mapper:      func(_ context.Context, n int) int { return n },
			want:        []int{7, 8, 9},
		},
		{
			name:        "single worker",
			jobs:        []int{3, 1, 4},
			workerCount: 1,
			mapper:      func(_ context.Context, n int) int { return n + 10 },
			want:        []int{13, 11, 14},
		},
		{
			name:        "empty input",
			jobs:        []int{},
			workerCount: 2,
			mapper:      func(_ context.Context, n int) int { return n },
			want:        []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MapConcurrent(tc.jobs, tc.workerCount, tc.mapper)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMapConcurrent_PreservesOrder(t *testing.T) {
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i
	}

	got := MapConcurrent(jobs, 8, func(_ context.Context, n int) int { return n * n })

	for i, v := range got {
		assert.Equal(t, i*i, v, "index %d", i)
	}
}

func TestMapConcurrent_StringTransform(t *testing.T) {
	tests := []struct {
		name        string
		jobs        []string
		workerCount int
		want        []string
	}{
		{"upper", []string{"hello", "world"}, 2, []string{"HELLO", "WORLD"}},
		{"single element", []string{"go"}, 1, []string{"GO"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MapConcurrent(tc.jobs, tc.workerCount, func(_ context.Context, s string) string {
				result := make([]byte, len(s))
				for i := range s {
					if s[i] >= 'a' && s[i] <= 'z' {
						result[i] = s[i] - 32
					} else {
						result[i] = s[i]
					}
				}
				return string(result)
			})
			assert.Equal(t, tc.want, got)
		})
	}
}
