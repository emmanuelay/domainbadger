package app

import (
	"errors"

	"github.com/emmanuelay/badger/pkg/whois"
	whoisparser "github.com/likexian/whois-parser-go"
)

func lookupDomain(domain string) DomainLookupResult {

	response, lookupErr := whois.Lookup(domain)

	body := string(response)
	result, whoisErr := whoisparser.Parse(body)

	lookupResult := DomainLookupResult{
		Domain:    domain,
		Available: errors.Is(whoisErr, whoisparser.ErrDomainNotFound) || whoisparser.IsNotFound(body),
		Error:     lookupErr,
		WhoIs:     result,
	}

	return lookupResult
}
