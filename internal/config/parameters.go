package config

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/zonedb/zonedb"
)

// Configuration ...
type Configuration struct {
	AllCharacters  bool
	Alpha          bool
	AlphaNumeric   bool
	Numeric        bool
	CustomRange    string
	Delay          int64
	TLD            []string
	SearchPatterns []string
}

func isValidDomain(domain string) bool {
	// This only validates the domain name, not the tld
	// Checks for alphanumeric characters including dash (-)
	// no spaces, min length of 1 and a max length of 63

	// RegEx explanation:
	// - 1st char only alphanumeric
	// - Subsequent chars (min 0, max 61) alphanumeric and dash
	// - Last char only alphanumeric
	r := regexp.MustCompile("(^[a-z0-9])([a-z0-9-]{0,61})(?:[a-z0-9])$")
	return r.MatchString(domain)
}

func isValidTLD(tld string) bool {
	r := regexp.MustCompile("(^[a-z]{2,16})$")
	return r.MatchString(tld)
}

func isValidRange(customRange string) bool {
	r := regexp.MustCompile("([a-z0-9-]{1,60})$")
	return r.MatchString(customRange)
}

func countWildcards(search string) int {
	wildcardFind := regexp.MustCompile("\\_")
	matches := wildcardFind.FindAllStringIndex(search, -1)
	return len(matches)
}

func validateArguments(config Configuration) error {

	if config.AllCharacters {
		config.Alpha = false
		config.AlphaNumeric = false
	}

	if config.AlphaNumeric == true {
		config.Alpha = false
	}

	if len(config.CustomRange) > 0 {
		config.Alpha = false
		config.AlphaNumeric = false
		config.AllCharacters = false

		// Check custom range for invalid characters
		if isValidRange(config.CustomRange) != true {
			return fmt.Errorf("Invalid custom characters specified: '%v'", config.CustomRange)
		}
	}

	// Make sure TLDs have a corresponding nameserver
	for _, tld := range config.TLD {

		if !isValidTLD(tld) {
			return fmt.Errorf("Invalid format TLD: '%v'", tld)
		}

		if zone := zonedb.PublicZone(tld); zone == nil {
			return fmt.Errorf("Invalid TLD specified: '%v'", tld)
		}
	}

	if len(config.SearchPatterns) == 0 {
		return errors.New("No search patterns provided")
	}

	// Check searchpatterns for wildcard character (underscore)
	// and invalid characters
	for _, search := range config.SearchPatterns {
		// if count := countWildcards(search); count == 0 {
		// return fmt.Errorf("Invalid search pattern, no wildcards found: '%v", search)
		// }

		clean := strings.ReplaceAll(search, "_", "")
		if len(clean) > 1 && !isValidDomain(clean) {
			return fmt.Errorf("Invalid search pattern, invalid domain: '%v", search)
		}
	}

	return nil
}

// GetConfigurationFromArguments ...
func GetConfigurationFromArguments() (Configuration, error) {

	config := Configuration{}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Println("\tbadger [parameters] <searchterms>")
		fmt.Println("\n\t<searchterms> are expected to use underscore as wildcard")
		fmt.Println("\nParameters:")

		flag.PrintDefaults()
	}

	flag.BoolVar(&config.AllCharacters, "all", true, "Use all possible characters (a-z, 0-9, -)")
	flag.BoolVar(&config.Alpha, "alpha", false, "Use alphabetic range (a-z)")
	flag.BoolVar(&config.AlphaNumeric, "alphanum", false, "Use alphanumeric range (a-z, 0-9)")
	flag.BoolVar(&config.Numeric, "numeric", false, "Use numeric range (0-9)")
	flag.StringVar(&config.CustomRange, "custom", "", "Use a custom character range (ex. abc123)")
	flag.Int64Var(&config.Delay, "delay", 500, "Delay between lookup attempts, in milliseconds")

	var tlds string
	flag.StringVar(&tlds, "tld", "com", "TLDs to search. Defaults to 'com'. Use comma to add multiple (ex. com,org,net).")

	flag.Parse()

	config.TLD = strings.Split(tlds, ",")

	config.SearchPatterns = flag.Args() // Search mask to use (ex. 'se_rchm_sk' to use 2 wildcard ranges)

	return config, validateArguments(config)
}
