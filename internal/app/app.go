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
	TLD          string
	TotalCount   int
	CurrentCount int
	Results      []lookupResult
	Bar          *mpb.Bar
}

// Run ...
func Run(cfg config.Configuration) {

	progressChannel := make(chan lookupResult, 1)
	doneChannel := make(chan string, 1)
	intDelay := int(cfg.Delay)
	alphaSet := []rune(cfg.Characters)
	tldCount := len(cfg.TLD)

	reports := map[string]report{}

	fmt.Println("Performing lookups")

	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(40))

	// For each TLD
	for _, tld := range cfg.TLD {

		rep := report{
			TLD:          tld,
			CurrentCount: 0,
			Results:      []lookupResult{},
		}

		// For each search pattern
		for _, searchPattern := range cfg.SearchPatterns {

			// Run generation of unique combinations
			uniqueCombinations := combinations.GenerateNames(alphaSet, searchPattern, "_")

			rep.TotalCount = rep.TotalCount + len(uniqueCombinations)

			wg.Add(1)
			// Run lookup for the TLD and the unique combinations generated from the search pattern
			go lookupDomainsForTLD(searchPattern, uniqueCombinations, tld, intDelay, progressChannel, doneChannel)
		}

		name := fmt.Sprintf("%v:", tld)
		bar := p.AddBar(int64(rep.TotalCount),
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

	for {
		select {
		case _ = <-doneChannel:
			{
				tldCount--
				wg.Done()
			}
		case result := <-progressChannel:
			{
				rep := reports[result.TLD]
				rep.CurrentCount++
				rep.Results = append(rep.Results, result)
				reports[result.TLD] = rep

				rep.Bar.Increment()
			}
		}

		if tldCount == 0 {
			break
		}
	}

	p.Wait()

	// TODO(ea): compile and display results nicely
	for idx, rep := range reports {
		fmt.Println("-", rep.TLD, len(rep.Results), idx, rep.CurrentCount, rep.TotalCount)
	}

}
