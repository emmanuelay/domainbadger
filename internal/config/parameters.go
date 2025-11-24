package config

import (
	"errors"
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

// ValidateConfiguration validates the provided configuration
func ValidateConfiguration(config Configuration) error {

	// If custom range is specified, it takes priority
	if len(config.CustomRange) > 0 {
		// Check custom range for invalid characters
		if !isValidRange(config.CustomRange) {
			return fmt.Errorf("invalid custom characters specified: '%v'", config.CustomRange)
		}
	}

	// Priority order: custom > numeric > alphanum > alpha > all
	// If user explicitly sets a specific flag, they don't want "all"

	// Make sure TLDs have a corresponding nameserver
	for _, tld := range config.TLD {

		if !isValidTLD(tld) {
			return fmt.Errorf("invalid format TLD: '%v'", tld)
		}

		if zone := zonedb.PublicZone(tld); zone == nil {
			return fmt.Errorf("invalid TLD specified: '%v'", tld)
		}
	}

	if len(config.SearchPatterns) == 0 {
		return errors.New("no search patterns provided")
	}

	// Check searchpatterns for wildcard character (underscore)
	// and invalid characters
	for _, search := range config.SearchPatterns {
		clean := strings.ReplaceAll(search, "_", "")
		if len(clean) > 1 && !isValidDomain(clean) {
			return fmt.Errorf("invalid search pattern, invalid domain: '%v", search)
		}
	}

	return nil
}

// GetCharacters returns the character range based on configuration flags
func (c *Configuration) GetCharacters() string {
	const (
		alphabet = "abcdefghijklmnopqrstuvwxyz"
		numerals = "0123456789"
		hyphen   = "-"
		all      = alphabet + numerals + hyphen
	)

	// Priority order: custom > numeric > alphanum > alpha > all (default)
	if len(c.CustomRange) > 0 {
		return c.CustomRange
	}

	if c.Numeric {
		return numerals
	}

	if c.AlphaNumeric {
		return alphabet + numerals
	}

	if c.Alpha {
		return alphabet
	}

	// Default to all characters (when no specific flag is set, or --all is specified)
	return all
}
