package utils

import "strings"

// LengthBetween returns true if min <= len(s) < max
func LengthBetween(s string, min, max int) bool {
	return len(s) >= min && len(s) < max
}

// ContainsAtLeast returns true if s contains at least num chars from chars
func ContainsAtLeast(s string, num int, chars string) bool {
	cnt := 0
	for _, c := range strings.Split(chars, "") {
		cnt += strings.Count(s, c)
	}
	return cnt >= num
}
