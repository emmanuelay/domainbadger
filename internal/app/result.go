package app

import whoisparser "github.com/likexian/whois-parser"

// Structs for reporting and accumulating progress (# available, # not available)

// DomainLookupResult encapsulates a single domain lookup
type DomainLookupResult struct {
	Domain    string
	TLD       string
	Available bool
	Error     error
	WhoIs     whoisparser.WhoisInfo
}

type DomainLookupJob struct {
	Domain string
	Delay  int64
}
