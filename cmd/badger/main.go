package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/emmanuelay/domainsearch/pkg/combinations"
	"github.com/emmanuelay/domainsearch/pkg/whois"

	whoisparser "github.com/likexian/whois-parser-go"
)

func domains(combos []string, tlds []string) []string {
	out := []string{}
	for _, tld := range tlds {
		for _, combo := range combos {
			domain := fmt.Sprintf("%v.%v", combo, tld)
			out = append(out, domain)
		}
	}
	return out
}

func main() {

	// TODO(ea): Move tlds to cli flag
	// TODO(ea): Move character range to cli flag
	// TODO(ea): Move search string to cli flag

	tlds := []string{"nu", "se"}
	alpha := "aoueiyåäö"
	inputSearch := "a*"
	starFind := regexp.MustCompile("\\*")
	matches := starFind.FindAllStringIndex(inputSearch, -1)
	fmt.Println(len(matches), "wildcards used")
	wildcardCount := len(matches)

	search := strings.ReplaceAll(inputSearch, "*", "%v")

	alphaSet := []rune(alpha)
	uniqueCombinations := combinations.Generate(alphaSet, search, wildcardCount)
	fmt.Println(len(uniqueCombinations), "domain name combinations generated")

	domains := domains(uniqueCombinations, tlds)
	fmt.Println(len(domains), "url combinations generated")

	for _, domain := range domains {

		time.Sleep(500 * time.Millisecond)

		response, err := whois.Lookup(domain)

		if err != nil {
			fmt.Println("Query Domain:", err.Error())
			return
		}
		body := string(response)

		result, err := whoisparser.Parse(body)

		if err != nil {

			if err == whoisparser.ErrDomainNotFound {
				fmt.Println(domain, "Domain not registered")
				continue
			}

			fmt.Println(err.Error())
			continue
		}

		if len(result.Domain.ExpirationDate) != 0 {
			fmt.Println(domain, "Domain expires at", result.Domain.ExpirationDate)
		}
	}
}
