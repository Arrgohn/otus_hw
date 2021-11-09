package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	resStream := in

	for _, stage := range stages {
		if stage == nil {
			continue
		}

		resStream = checkDone(done, stage(resStream))
	}

	return resStream
}

func checkDone(done In, in In) Out {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case ch <- val:
				}
			}
		}
	}()
	return ch
}
