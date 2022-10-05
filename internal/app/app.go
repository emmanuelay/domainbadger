package app

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
)

// Run ...
func Run(ctx context.Context, cfg config.Configuration) {
	intDelay := int(cfg.Delay)
	alphaSet := []rune(cfg.Characters)
	tldCount := len(cfg.TLD)

	allCombinations := []string{}

	// Generate all combinations
	for _, searchPattern := range cfg.SearchPatterns {
		patternCombinations := combinations.GenerateNames(alphaSet, searchPattern, "_")
		allCombinations = append(allCombinations, patternCombinations...)
	}

	fmt.Printf("Generated %d unique combinations\n", len(allCombinations))
	totalLookups := len(allCombinations) * tldCount
	fmt.Printf("Performing %v separate lookups\n", totalLookups)

	var wg sync.WaitGroup

	maxProcs := runtime.GOMAXPROCS(0)
	resultChannel := make(chan DomainLookupResult, totalLookups)
	jobsChannel := make(chan DomainLookupJob, maxProcs)

	fmt.Printf("Starting %d workers ...\n", maxProcs)

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for w := 1; w <= maxProcs; w++ {
		go worker(cancelCtx, w+1, &wg, intDelay, jobsChannel, resultChannel)
	}

	bar := progressbar.Default(int64(totalLookups), "Queuing")

	// This loop starts all lookups.
	for _, tld := range cfg.TLD {
		for _, combination := range allCombinations {
			domain := fmt.Sprintf("%s.%s", combination, tld)
			job := DomainLookupJob{
				Domain: domain,
			}
			jobsChannel <- job
			bar.Add(1)
		}
	}

	fmt.Println("All jobs added, waiting")

	tab := tablewriter.NewWriter(os.Stdout)
	tab.SetHeader([]string{"domain", "available", "error", "expires at", "last renewed", "registrant"})

	for i := 0; i < totalLookups; i++ {

		result := <-resultChannel

		var expirationDate, updateDate, registrant string

		if result.WhoIs.Domain != nil {
			expirationDate = result.WhoIs.Domain.ExpirationDate
			updateDate = result.WhoIs.Domain.UpdatedDate
		}

		if result.WhoIs.Registrant != nil {
			registrant = result.WhoIs.Registrant.Name
		}

		if result.Available {
			tab.Append([]string{
				result.Domain,
				"✅",
				"",
				"-",
				"-",
				"-",
			})

		} else {
			if result.Error != nil {
				tab.Append([]string{
					result.Domain,
					"❌",
					result.Error.Error(),
					expirationDate,
					updateDate,
					registrant,
				})
			} else {
				tab.Append([]string{
					result.Domain,
					"🛑",
					"",
					expirationDate,
					updateDate,
					registrant,
				})

			}
		}

	}

	close(jobsChannel)
	close(resultChannel)

	wg.Wait()

	tab.Render()
}
