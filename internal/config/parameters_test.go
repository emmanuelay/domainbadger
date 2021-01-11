package config

import (
	"testing"
)

func TestIsValidRange(t *testing.T) {
	if isValidRange("abcd") != true {
		t.Error("Proper range should be accepted")
	}

	if isValidRange("abcd!") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidRange("abcd*_") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidRange("abc-def") != true {
		t.Error("Proper range should be accepted")
	}

}

func TestCountWildcards(t *testing.T) {

	if count := countWildcards("h_ll_"); count != 2 {
		t.Error("Expected 2, got", count)
	}

	if count := countWildcards("w_rld"); count != 1 {
		t.Error("Expected 1, got", count)
	}

	if count := countWildcards("d_min_ati_n"); count != 3 {
		t.Error("Expected 3, got", count)
	}
}

func TestIsValidTLD(t *testing.T) {

	if isValidTLD("se") != true {
		t.Error("Should validate correct tld")
	}

	if isValidTLD("sesesesesesesesese") != false {
		t.Error("Invalid length tld")
	}

	if isValidTLD("s") != false {
		t.Error("Invalid length, short tld")
	}

	if isValidTLD("s1") != false {
		t.Error("Invalid tld with numeric character")
	}

	if (isValidTLD("s!") != false) || (isValidTLD("!2") != false) {
		t.Error("Invalid tld with invalid character")
	}

}

func TestIsValidDomain(t *testing.T) {

	if isValidDomain("-domain") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain-") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain!") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain?") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain=") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain,") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain.") != false {
		t.Error("Invalid character should not be accepted")
	}

	if isValidDomain("domain") != true {
		t.Error("Domain should have been accepted")
	}

	if isValidDomain("domain-domain") != true {
		t.Error("Domain should have been accepted")
	}

	if isValidDomain("domain-domain-domain") != true {
		t.Error("Domain should have been accepted")
	}

	if isValidDomain("DoMaIn") != false {
		t.Error("Should not accept uppercase characters")
	}

	if isValidDomain("abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopq") != false {
		t.Error("Too long domain should not be accepted")
	}

	if isValidDomain("abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz") != true {
		t.Error("Maxed out domain name should be accepted")
	}

	if isValidDomain("x1") != true {
		t.Error("Minimal domain name should be accepted")
	}

	if isValidDomain("x-") != false {
		t.Error("Minimal domain name with invalid char should not be accepted")
	}
}
