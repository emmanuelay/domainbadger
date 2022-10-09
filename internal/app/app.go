package app

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
)

// Run ...
func Run(ctx context.Context, cfg config.Configuration) {

	domains := combinations.GenerateDomainCombinations(cfg.Characters, cfg.SearchPatterns, cfg.TLD)
	totalLookups := len(domains)

	fmt.Printf("Generated %d unique combinations\n", totalLookups)

	maxProcs := 20
	resultChannel := make(chan DomainLookupResult)
	jobsChannel := make(chan DomainLookupJob, maxProcs)

	quit := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(maxProcs)

	go func() {
		// Close the result channel when all workers have finished
		wg.Wait()
		close(resultChannel)
	}()

	fmt.Printf("Starting %d workers ...\n", maxProcs)
	for w := 0; w < maxProcs; w++ {
		go func(workerId int) {
			defer wg.Done()
			worker(ctx, workerId, jobsChannel, resultChannel)
		}(w)
	}

	bar := progressbar.Default(int64(totalLookups), "Querying")

	go func() {
		output := [][]string{}

		for result := range resultChannel {
			bar.Add(1)

			var expirationDate, updateDate, registrant string

			if result.WhoIs.Domain != nil {
				expirationDate = result.WhoIs.Domain.ExpirationDate
				updateDate = result.WhoIs.Domain.UpdatedDate
			}

			if result.WhoIs.Registrant != nil {
				registrant = result.WhoIs.Registrant.Name
			}

			if result.Available {
				output = append(output, []string{result.Domain, "✅", "", "-", "-", "-"})
			} else {
				if result.Error != nil {
					// Lookup failed
					output = append(output, []string{result.Domain, "❌", result.Error.Error(), expirationDate, updateDate, registrant})
				} else {
					// Domain is already registered
					output = append(output, []string{result.Domain, "🛑", "", expirationDate, updateDate, registrant})
				}
			}
		}
		if err := bar.Close(); err != nil {
			fmt.Println(err.Error())
		}

		tab := tablewriter.NewWriter(os.Stdout)
		tab.SetHeader([]string{"domain", "available", "error", "expires at", "last renewed", "registrant"})

		// TODO(ea): Sort the output

		tab.AppendBulk(output)
		tab.Render()

		quit <- true
	}()

	// Push all lookups to queue/jobs channel
	for _, domain := range domains {
		job := DomainLookupJob{
			Domain: domain,
			Delay:  cfg.Delay,
		}
		jobsChannel <- job

	}
	close(jobsChannel)

	<-quit
}
