package app

import (
	"fmt"
	"sync"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"

	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

type report struct {
	TLD     string
	Results []lookupResult
	Bar     *mpb.Bar
}

// Run ...
func Run(cfg config.Configuration) {

	progressChannel := make(chan lookupResult, 1)
	doneChannel := make(chan string, 1)
	intDelay := int(cfg.Delay)
	alphaSet := []rune(cfg.Characters)
	tldCount := len(cfg.TLD)

	reports := map[string]report{}
	allCombinations := []string{}

	// Generate all combinations
	for _, searchPattern := range cfg.SearchPatterns {
		patternCombinations := combinations.GenerateNames(alphaSet, searchPattern, "_")
		allCombinations = append(allCombinations, patternCombinations...)
	}
	fmt.Printf("Generated %d combinations\n", len(allCombinations))

	wg := sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(40))

	fmt.Printf("Performing %v lookups\n", len(allCombinations)*tldCount)

	// This loop starts all lookups.
	// Grouped by TLD because we want progress reported on a TLD-basis
	for _, tld := range cfg.TLD {

		rep := report{
			TLD:     tld,
			Results: []lookupResult{},
		}

		// Run lookup for the TLD and the unique combinations generated from the search pattern
		go lookupDomainsForTLD(&wg, allCombinations, tld, intDelay, progressChannel, doneChannel)

		// Add progress bar for current TLD
		name := fmt.Sprintf("%v:", tld)
		bar := p.AddBar(
			int64(len(allCombinations)),
			mpb.PrependDecorators(
				decor.Name(name, decor.WC{W: 6, C: decor.DidentRight}),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.AverageETA(decor.ET_STYLE_GO), " ",
				),
			),
		)
		rep.Bar = bar
		reports[tld] = rep
	}

	// This loops waits for all lookups to finish
	for {
		select {
		case _ = <-doneChannel:
			{
				tldCount--
			}
		case result := <-progressChannel:
			{
				rep := reports[result.TLD]
				rep.Results = append(rep.Results, result)
				reports[result.TLD] = rep

				rep.Bar.Increment()
			}
		}

		if tldCount == 0 {
			break
		}
	}

	// Wait for progress bar to render properly
	p.Wait()

	// TODO(ea): compile and display results nicely
	for idx, rep := range reports {
		fmt.Println("-", rep.TLD, len(rep.Results), idx)
	}

}
