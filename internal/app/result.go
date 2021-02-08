package app

import whoisparser "github.com/likexian/whois-parser-go"

// Structs for reporting and accumulating progress (# available, # not available)

// lookupResult encapsulates a single lookup
type lookupResult struct {
	Domain    string
	TLD       string
	Available bool
	Error     error
	WhoIs     whoisparser.WhoisInfo
}

// tldResult encapsulates a series of lookups for a tld
type tldResult struct {
	TLD     string
	Lookups []lookupResult
}
