package combinations

import (
	"fmt"
	"regexp"
	"strings"
)

// GenerateNames creates a set of combinations based on which characters are allowed (set)
func GenerateNames(set []rune, searchPattern string, wildcardCharacter string) []string {

	countWildcardCharacters := countWildcards(searchPattern) // "H_ll_"
	searchString := strings.ReplaceAll(searchPattern, wildcardCharacter, "%v")

	res := printAllKLengthRec(set, "", len(set), countWildcardCharacters)
	out := []string{}
	for _, v := range res {
		out = append(out, combine(searchString, v))
	}
	return out
}

// GenerateDomains creates a set of domains using a set of TLDs
func GenerateDomains(combos []string, tlds []string) []string {
	out := make([]string, 0, len(combos)*len(tlds))
	for _, tld := range tlds {
		for _, combo := range combos {
			domain := fmt.Sprintf("%v.%v", combo, tld)
			out = append(out, domain)
		}
	}
	return out
}

func GenerateDomainCombinations(characters string, searchPatterns, domains []string) []string {
	alphaSet := []rune(characters)
	allCombinations := []string{}

	// Generate all combinations
	for _, searchPattern := range searchPatterns {
		patternCombinations := GenerateNames(alphaSet, searchPattern, "_")
		allCombinations = append(allCombinations, patternCombinations...)
	}

	return GenerateDomains(allCombinations, domains)
}

func countWildcards(search string) int {
	wildcardFind := regexp.MustCompile(`\_`)
	matches := wildcardFind.FindAllStringIndex(search, -1)
	return len(matches)
}

func printAllKLengthRec(set []rune, prefix string, n int, k int) []string {

	if k == 0 {
		return []string{prefix}
	}

	out := []string{}

	for i := 0; i < n; i++ {
		newPrefix := prefix + string(set[i])
		res := printAllKLengthRec(set, newPrefix, n, k-1)
		out = append(out, res...)
	}
	return out
}

func combine(pattern string, combination string) string {
	comboRunes := []rune(combination)
	tmp := make([]interface{}, len(comboRunes))
	for i, val := range comboRunes {
		tmp[i] = string(val)
	}
	return fmt.Sprintf(pattern, tmp...)
}
