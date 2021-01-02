package main

import "testing"

func TestValidateDomain(t *testing.T) {

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

	if validateDomain("abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890") != false {
		t.Error("Too long domain should not be accepted")
	}

}
