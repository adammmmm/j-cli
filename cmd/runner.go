package cmd

import (
	"context"
	"sync"
)

type DeviceTask func(ctx context.Context, device string) workerResult
type workerResult struct {
	device string
	output string
	err    error
}

func runOnDevices(
	ctx context.Context,
	devices []string,
	workers int,
	task DeviceTask,
) <-chan workerResult {

	results := make(chan workerResult)
	jobs := make(chan string)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for device := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					results <- task(ctx, device)
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, d := range devices {
			select {
			case <-ctx.Done():
				return
			case jobs <- d:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
