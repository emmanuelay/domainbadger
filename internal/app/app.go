package app

import (
	"fmt"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
)

// Run ...
func Run(cfg config.Configuration) {

	progressChannel := make(chan lookupResult, 1)
	doneChannel := make(chan string, 1)
	intDelay := int(cfg.Delay)
	alphaSet := []rune(cfg.Characters)
	tldCount := 0

	// For each search pattern
	for _, searchPattern := range cfg.SearchPatterns {

		// Run generation of unique combinations
		uniqueCombinations := combinations.GenerateNames(alphaSet, searchPattern, "_")

		// For each TLD
		for _, tld := range cfg.TLD {

			// Run lookup for the TLD and the unique combinations generated from the search pattern
			go lookupDomainsForTLD(searchPattern, uniqueCombinations, tld, intDelay, progressChannel, doneChannel)
		}

		tldCount = tldCount + len(cfg.TLD)
	}

	results := map[string][]lookupResult{}

	// TODO(ea): Create Multi Progress Bar
	for {
		select {
		case val := <-doneChannel:
			{
				// TODO(ea): Update progress bar
				fmt.Println("Scan completed for TLD:", val)
				tldCount--
			}
		case result := <-progressChannel:
			{
				// TODO(ea): Update progress bar
				if result.Available {
					fmt.Println("-", result.Domain, "is available")
				}

				results[result.TLD] = append(results[result.TLD], result)
			}
		}

		if tldCount == 0 {
			break
		}
	}

	// TODO(ea): compile and display results nicely

	for idx, tld := range results {
		fmt.Println("-", len(tld), idx)
	}

}
