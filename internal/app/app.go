package app

import (
	"fmt"
	"time"

	"github.com/emmanuelay/badger/internal/config"
	"github.com/emmanuelay/badger/pkg/combinations"
	"github.com/emmanuelay/badger/pkg/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

func getCharacterRange(customRange string, alphaNum, num bool) string {
	const (
		alphabet = "abcdefghijklmnopqrstuvwxyz"
		numerals = "0123456789"
		all      = alphabet + numerals
	)

	if len(customRange) > 0 {
		return customRange
	}

	if alphaNum {
		return all
	}
	if num {
		return numerals
	}

	// Default to alpha
	return alphabet
}

// Run ...
func Run(cfg config.Configuration) {

	alpha := getCharacterRange(cfg.CustomRange, cfg.AlphaNumeric, cfg.Numeric)

	count := 0

	// TODO(ea): Create struct for reporting progress (current, # registered, # not registered etc) in a structured way
	doneChannel := make(chan string, 1)

	for _, search := range cfg.SearchPatterns {
		fmt.Println("Performing generation for:", search)

		alphaSet := []rune(alpha)

		// Run generation of unique combinations
		uniqueCombinations := combinations.GenerateNames(alphaSet, search, "_")
		fmt.Println(len(uniqueCombinations), "domain name combinations generated")

		// Parallellize lookup for each TLD
		for _, tld := range cfg.TLD {
			go lookupDomainsForTLD(search, uniqueCombinations, tld, int(cfg.Delay), doneChannel)
		}

		count = count + len(cfg.TLD)
	}

	// TODO(ea): Add Multi Progress Bar
	for {
		select {
		case val := <-doneChannel:
			fmt.Println("Scan completed for TLD:", val)
			count--
		}

		if count == 0 {
			break
		}
	}

	// TODO(ea): compile and display results nicely
}

func lookupDomainsForTLD(pattern string, names []string, tld string, delay int, done chan string) {
	domains := combinations.GenerateDomains(names, []string{tld})
	for _, domain := range domains {
		lookupDomain(domain)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	done <- fmt.Sprintf("%v.%v", pattern, tld)
}

func lookupDomain(domain string) {

	response, err := whois.Lookup(domain)

	if err != nil {
		fmt.Println("Query Domain:", err.Error())
		return
	}

	body := string(response)

	result, err := whoisparser.Parse(body)

	if err == whoisparser.ErrDomainNotFound || whoisparser.IsNotFound(body) {
		fmt.Println(domain, "Domain not registered")
		return
	}

	if err != nil {
		fmt.Println(domain, err.Error())
		return
	}

	if result.Domain != nil {
		if len(result.Domain.ExpirationDate) > 0 {
			fmt.Println(domain, "Domain expires at", result.Domain.ExpirationDate)
			return
		}

		if len(result.Domain.UpdatedDate) > 0 {
			fmt.Println(domain, "Domain last updated at", result.Domain.UpdatedDate)
			return
		}
	}
}
