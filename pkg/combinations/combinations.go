package combinations

import "fmt"

// Generate creates a set of combinations based on which characters are allowed (set)
func Generate(set []rune, searchString string, wildcardCharacters int) []string {
	res := printAllKLengthRec(set, "", len(set), wildcardCharacters)
	out := []string{}
	for _, v := range res {
		out = append(out, combine(searchString, v))
	}
	return out
}

// The main recursive method to print all possible strings of length k
func printAllKLengthRec(set []rune, prefix string, n int, k int) []string {

	// Base case: k is 0,
	// print prefix
	if k == 0 {
		return []string{prefix}
	}

	// One by one add all characters
	// from set and recursively
	// call for k equals to k-1

	out := []string{}

	for i := 0; i < n; i++ {

		// Next character of input added
		newPrefix := prefix + string(set[i])

		// k is decreased, because
		// we have added a new character
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
