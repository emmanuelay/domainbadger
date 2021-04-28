package app

import whoisparser "github.com/likexian/whois-parser-go"

// Structs for reporting and accumulating progress (# available, # not available)

// TLDResults encapsulates all lookups performed within a TLD
type TLDResults struct {
	TLD         string
	Results     []DomainLookupResult
	TotalCount  int64
	Available   int64
	Unavailable int64
	ErrorCount  int64
}

// DomainLookupResult encapsulates a single domain lookup
type DomainLookupResult struct {
	Domain    string
	TLD       string
	Available bool
	Error     error
	WhoIs     whoisparser.WhoisInfo
}
