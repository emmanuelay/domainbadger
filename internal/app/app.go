package app

import (
	"fmt"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
)

// Run ...
func Run(cfg config.Configuration) {

	tldCount := 0

	progressChannel := make(chan lookupResult, 1)
	doneChannel := make(chan string, 1)

	for _, search := range cfg.SearchPatterns {
		fmt.Println("Performing generation for:", search)

		alphaSet := []rune(cfg.Characters)

		// Run generation of unique combinations
		uniqueCombinations := combinations.GenerateNames(alphaSet, search, "_")

		// Parallellize lookup for each TLD
		for _, tld := range cfg.TLD {
			go lookupDomainsForTLD(search, uniqueCombinations, tld, int(cfg.Delay), progressChannel, doneChannel)
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
