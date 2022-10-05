package app

import (
	"context"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup, delay int, jobs <-chan DomainLookupJob, results chan<- DomainLookupResult) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			results <- lookupDomain(job.Domain)
			time.Sleep(time.Duration(delay))

		case <-ctx.Done():
			return
		}
	}
}
