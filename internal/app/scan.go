package app

// Structs for reporting and accumulating progress (# available, # not available)

// lookupResult encapsulates a single lookup
type lookupResult struct {
	Domain    string
	Available bool
}

// tldResult encapsulates a series of lookups for a tld
type tldResult struct {
	TLD     string
	Lookups []lookupResult
}
