package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
)

// Run ...
func Run(cfg config.Configuration) {

	progressChannel := make(chan DomainLookupResult, 1)
	doneChannel := make(chan string, 1)
	intDelay := int(cfg.Delay)
	alphaSet := []rune(cfg.Characters)
	tldCount := len(cfg.TLD)

	reports := map[string]TLDResults{}
	allCombinations := []string{}

	// Generate all combinations
	for _, searchPattern := range cfg.SearchPatterns {
		patternCombinations := combinations.GenerateNames(alphaSet, searchPattern, "_")
		allCombinations = append(allCombinations, patternCombinations...)
	}

	fmt.Printf("Generated %d unique combinations\n", len(allCombinations))

	wg := sync.WaitGroup{}

	fmt.Printf("Performing %v separate lookups\n", len(allCombinations)*tldCount)

	// This loop starts all lookups.
	// Grouped by TLD because we want progress reported on a TLD-basis
	for _, tld := range cfg.TLD {

		rep := TLDResults{
			TLD:         tld,
			Results:     []DomainLookupResult{},
			Available:   0,
			Unavailable: 0,
			ErrorCount:  0,
			TotalCount:  0,
		}

		// Run lookup for the TLD and the unique combinations generated from the search pattern
		go lookupDomainsForTLD(&wg, allCombinations, tld, intDelay, progressChannel, doneChannel)

		// Add progress bar for current TLD
		reports[tld] = rep
	}

	// This loops receives results from lookups
	// ..and breaks when the last lookup is done
	for {
		select {
		case tldTime := <-doneChannel:
			{
				fmt.Println(tldTime)
				tldCount--
			}
		case result := <-progressChannel:
			{
				rep := reports[result.TLD]

				rep.TotalCount++

				if result.Available {
					fmt.Println("✅", result.Domain)
					rep.Available++
				} else {
					if result.Error != nil {
						fmt.Println("❌", result.Domain, result.Error)
						rep.ErrorCount++
					} else {
						fmt.Println("🛑", result.Domain)
						rep.Unavailable++
					}
				}

				rep.Results = append(rep.Results, result)
				reports[result.TLD] = rep
			}
		}

		if tldCount == 0 {
			break
		}
	}

	// Display summary
	fmt.Println("-")

	for _, rep := range reports {

		res := fmt.Sprintf("[.%s]\t%v out of %v domains available", strings.ToUpper(rep.TLD), rep.Available, rep.TotalCount)

		if rep.Unavailable > 0 {
			res += fmt.Sprintf(", %v unavailable", rep.Unavailable)
		}

		if rep.ErrorCount > 0 {
			res += fmt.Sprintf(", %v failed lookup", rep.ErrorCount)
		}

		fmt.Println(res)
	}

}
