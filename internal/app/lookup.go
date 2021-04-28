package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/emmanuelay/badger/pkg/combinations"
	"github.com/emmanuelay/badger/pkg/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

func lookupDomainsForTLD(waitGroup *sync.WaitGroup, names []string, tld string, delay int, progress chan DomainLookupResult, done chan string) {

	waitGroup.Add(1)
	start := time.Now()

	domains := combinations.GenerateDomains(names, []string{tld})
	for _, domain := range domains {

		lookup := lookupDomain(domain)
		lookup.TLD = tld

		if lookup.Error != nil {
			if lookup.Error.Error() == "Domain query limit exceeded." {
				fmt.Printf("Query limit reached for '.%v' - ", tld)
				break
			}
		}

		// Pass results to main thread
		progress <- lookup

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	duration := time.Since(start)

	done <- fmt.Sprintf(".%v finished in (%.2f seconds)", tld, duration.Seconds())
	waitGroup.Done()
}

func lookupDomain(domain string) DomainLookupResult {

	response, err := whois.Lookup(domain)

	lookupResult := DomainLookupResult{
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
