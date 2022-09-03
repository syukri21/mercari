package workers

import (
	"context"
)

// DispatchWorkers ...
type DispatchWorkers[Data any] struct {
	TotalWorker int
	Ctx         context.Context
}

func NewDispatchWorkers[Data any](totalWorker int, ctx context.Context) *DispatchWorkers[Data] {
	return &DispatchWorkers[Data]{TotalWorker: totalWorker, Ctx: ctx}
}

func (u *DispatchWorkers[Data]) DispatchJobs(dispatchFunc func(ctx context.Context, insertData Data) error, onDone func(counter int, indexWorker int), jobs <-chan Data, errChan chan<- error) {
	for index := 0; index < u.TotalWorker; index++ {
		go func(
			index int,
			jobs <-chan Data,
			ctx context.Context,
			onDone func(counter int, indexWorker int),
		) {
			counter := 0
			for job := range jobs {
				err := dispatchFunc(ctx, job)
				if err != nil {
					errChan <- err
					return
				}
				onDone(counter, index)
				counter++
			}
		}(index, jobs, u.Ctx, onDone)
	}
}
