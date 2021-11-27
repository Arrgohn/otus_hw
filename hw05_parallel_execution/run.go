package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		m = 0
	}

	errs := make(chan error, m)
	ch := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}

	wg.Add(n)

	for _, task := range tasks {
		ch <- task
	}

	for i := 0; i < n; i++ {
		go worker(ch, errs, &wg)
	}

	wg.Wait()

	err := error(nil)
	select {
	case errs <- err:
		return nil
	default:
		if m > 0 {
			return ErrErrorsLimitExceeded
		}
		return nil
	}
}

func worker(tasks chan Task, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task := <-tasks:
			err := task()
			if err != nil {
				select {
				case errs <- err:
					continue
				default:
					return
				}
			}
		default:
			return
		}
	}
}
