package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/emmanuelay/domainsearch/internal/config"
	"github.com/emmanuelay/domainsearch/pkg/combinations"
	"github.com/emmanuelay/domainsearch/pkg/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

func domains(combos []string, tlds []string) []string {
	out := []string{}
	for _, tld := range tlds {
		for _, combo := range combos {
			domain := fmt.Sprintf("%v.%v", combo, tld)
			out = append(out, domain)
		}
	}
	return out
}

func countWildcards(search string) int {
	wildcardFind := regexp.MustCompile("\\_")
	matches := wildcardFind.FindAllStringIndex(search, -1)
	return len(matches)
}

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	numerals = "0123456789"
	all      = alphabet + numerals
)

// Run ...
func Run(cfg config.Configuration) {

	alpha := ""

	if len(cfg.CustomRange) > 0 {
		alpha = cfg.CustomRange
	} else {
		if cfg.Alpha {
			alpha = alphabet
		}
		if cfg.AlphaNumeric {
			alpha = all
		}
		if cfg.Numeric {
			alpha = numerals
		}
	}

	for _, search := range cfg.SearchPatterns {
		fmt.Println("Performing generation for:", search)
		count := countWildcards(search)
		searchPattern := strings.ReplaceAll(search, "_", "%v")

		alphaSet := []rune(alpha)

		// Run generation of unique combinations
		uniqueCombinations := combinations.Generate(alphaSet, searchPattern, count)
		fmt.Println(len(uniqueCombinations), "domain name combinations generated")

		// Run generation of domains to check
		domains := domains(uniqueCombinations, cfg.TLD)
		fmt.Println(len(domains), "url combinations generated")

		// run whoislookups
		for _, domain := range domains {

			lookupDomain(domain)

			time.Sleep(time.Duration(cfg.Delay) * time.Millisecond)
		}

		// TODO(ea): compile and display results nicely
		fmt.Println(" ")
	}

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
