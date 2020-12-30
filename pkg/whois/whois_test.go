package whois

import (
	"testing"
)

func TestWhois(t *testing.T) {

	result, err := Lookup("google.com")
	if err != nil {
		t.Errorf("Error while performing whois (%v)", err)
	}

	if len(result) != 3559 {
		t.Errorf("Expected 3559 bytes response, got %v", len(result))
	}
}
