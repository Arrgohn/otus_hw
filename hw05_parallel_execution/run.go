package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error
var wg sync.WaitGroup
var errs chan error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errs = make(chan error, m)
	ch := make(chan Task, len(tasks))
	wg = sync.WaitGroup{}
	wg.Add(n)

	for _, task := range tasks {
		ch <- task
	}

	for i := 0; i < n; i++ {
		go worker(ch)
	}

	wg.Wait()

	err := error(nil)
	select {
		case errs <- err:
			return nil
		default:
			return ErrErrorsLimitExceeded
		}
}

func worker(tasks chan Task)  {
	for {
		select {
		case task := <-tasks:
			err := task()

			if err != nil {
				select {
				case errs <- err:
					continue
				default:
					wg.Done()
					return
				}
			}
		default:
			wg.Done()
			return
		}
	}
}

