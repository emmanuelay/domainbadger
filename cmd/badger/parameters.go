package main

import (
	"flag"
)

type Configuration struct {
	Alpha          bool
	AlphaNumeric   bool
	CustomRange    string
	TLD            []string
	SearchPatterns []string
}

func getConfigurationFromArguments() Configuration {

	config := Configuration{}

	flag.BoolVar(&config.Alpha, "alpha", true, "Use alphabetic range (a-z)")
	flag.BoolVar(&config.AlphaNumeric, "alphanum", false, "Use alphanumeric range (a-z, 0-9)")
	flag.StringVar(&config.CustomRange, "custom", "", "Use a custom character range (ex. abc123)")

	var tlds string
	flag.StringVar(&tlds, "tld", "com", "TLDs to search. Use comma to add multiple (ex. com,org,net)")

	flag.Parse()

	if config.AlphaNumeric == true {
		config.Alpha = false
	}

	if len(config.CustomRange) > 0 {
		config.Alpha = false
		config.AlphaNumeric = false
	}

	config.SearchPatterns = flag.Args() // Search mask to use (ex. 'se*rchm*sk' to use 2 wildcard ranges)

	return config
}
