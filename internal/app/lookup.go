package app

import (
	"fmt"
	"time"

	"github.com/emmanuelay/badger/pkg/combinations"
	"github.com/emmanuelay/badger/pkg/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

func lookupDomainsForTLD(pattern string, names []string, tld string, delay int, progress chan lookupResult, done chan string) {

	start := time.Now()

	domains := combinations.GenerateDomains(names, []string{tld})
	for _, domain := range domains {

		// TODO(ea): pass results to main thread
		lookup := lookupDomain(domain)
		lookup.TLD = tld
		progress <- lookup

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	duration := time.Since(start)

	done <- fmt.Sprintf("%v.%v (%v s)", pattern, tld, duration.Seconds())
}

func lookupDomain(domain string) lookupResult {

	response, err := whois.Lookup(domain)

	lookupResult := lookupResult{
		Domain:    domain,
		Available: false,
	}

	if err != nil {
		fmt.Println("Query Domain:", err.Error())
		lookupResult.Error = err
		return lookupResult
	}

	body := string(response)

	result, err := whoisparser.Parse(body)

	lookupResult.WhoIs = result

	if err == whoisparser.ErrDomainNotFound || whoisparser.IsNotFound(body) {
		//fmt.Println(domain, "Domain not registered")
		lookupResult.Available = true
		return lookupResult
	}

	if err != nil {
		fmt.Println(domain, err.Error())
		lookupResult.Error = err
		return lookupResult
	}

	if result.Domain != nil {
		if len(result.Domain.ExpirationDate) > 0 {
			//fmt.Println(domain, "Domain expires at", result.Domain.ExpirationDate)
			return lookupResult
		}

		if len(result.Domain.UpdatedDate) > 0 {
			//fmt.Println(domain, "Domain last updated at", result.Domain.UpdatedDate)
			return lookupResult
		}
	}

	return lookupResult
}
