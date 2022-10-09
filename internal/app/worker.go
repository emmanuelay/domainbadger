package app

import (
	"context"
	"time"
)

func worker(ctx context.Context, id int, jobs <-chan DomainLookupJob, results chan<- DomainLookupResult) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			result := lookupDomain(job.Domain)
			results <- result
			time.Sleep(time.Duration(job.Delay))

		case <-ctx.Done():
			return
		}
	}
}
