package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("if m=0 then error limit exceeded", func(t *testing.T) {
		n := 10
		m := 0

		task := func() error {
			return nil
		}
		tasks := make([]Task, 0, 50)
		for i := 0; i < cap(tasks); i++ {
			tasks = append(tasks, task)
		}
		err := Run(tasks, n, m)

		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("test concurrency without time.Sleep", func(t *testing.T) {
		// n = 1. Get time (t1) after running 5000 simple tasks
		n := 1
		task := func() error {
			for a := 1; a < 100; a++ {
				x := rand.Intn(a)
				x *= rand.Intn(a)
				x /= rand.Intn(a) + 1

				if x != x-1 { // Workaround for linter which doesn't allow last x result to be unused
					return nil
				}
			}
			return nil
		}
		tasks := make([]Task, 0, 5000)
		for i := 0; i < cap(tasks); i++ {
			tasks = append(tasks, task)
		}

		start := time.Now()
		err := Run(tasks, n, 1)
		t1 := time.Since(start)

		require.NoError(t, err)

		// n = 2. Test t2 less than t1 * 2
		n = 2
		start = time.Now()
		err = Run(tasks, n, 1)
		t2 := time.Since(start)

		require.NoError(t, err)
		require.Less(t, t2, t1*2)

		// n = 10. Test t3 less than t1 * 10
		n = 10
		start = time.Now()
		err = Run(tasks, n, 1)
		t3 := time.Since(start)

		require.NoError(t, err)
		require.Less(t, t3, t1*10)

		// n = 5000. Test t5000 less than t1 * 10
		n = 5000
		start = time.Now()
		err = Run(tasks, n, 1)
		t5000 := time.Since(start)

		require.NoError(t, err)
		require.Less(t, t5000, t1*10) // compare with t1*10 (just 10, not 5000)
	})
}
