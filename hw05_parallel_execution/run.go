package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	errCount := int32(0)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			worker(&tasks, m, &errCount)
		}()
	}

	wg.Wait()

	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(tasks *[]Task, maxErr int, errCount *int32) {
	mu := sync.Mutex{}
	for {
		mu.Lock()

		if len(*tasks) == 0 {
			mu.Unlock()
			return
		}

		t := (*tasks)[0]
		*tasks = (*tasks)[1:]
		mu.Unlock()

		err := t()

		if err != nil {
			atomic.AddInt32(errCount, 1)
		}

		if *errCount >= int32(maxErr) {
			return
		}
	}
}
