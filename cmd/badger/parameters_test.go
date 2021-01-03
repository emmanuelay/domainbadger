package main

import "testing"

func TestValidateDomain(t *testing.T) {

	if validateDomain("-domain") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain-") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain!") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain?") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain=") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain,") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain.") != false {
		t.Error("Invalid character should not be accepted")
	}

	if validateDomain("domain") != true {
		t.Error("Domain should have been accepted")
	}

	if validateDomain("domain-domain") != true {
		t.Error("Domain should have been accepted")
	}

	if validateDomain("domain-domain-domain") != true {
		t.Error("Domain should have been accepted")
	}

	if validateDomain("DoMaIn") != false {
		t.Error("Should not accept uppercase characters")
	}

	if validateDomain("abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopq") != false {
		t.Error("Too long domain should not be accepted")
	}

	if validateDomain("abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz") != true {
		t.Error("Maxed out domain name should be accepted")
	}

	if validateDomain("x1") != true {
		t.Error("Minimal domain name should be accepted")
	}

	if validateDomain("x-") != false {
		t.Error("Minimal domain name with invalid char should not be accepted")
	}
}
