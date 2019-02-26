package utils

import "regexp"

// Use regular expression to match the IPv4 addresses.
// Ref: https://www.hackerrank.com/challenges/java-regex/forum
func IsIP(str string) bool {
	regex := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	return regex.MatchString(str)
}

// Use regular expression to match the domain name.
// Ref: https://stackoverflow.com/questions/10306690/what-is-a-regular-expression-which-will-match-a-valid-domain-name-without-a-subd
func IsDomainName(str string) bool {
	regex := regexp.MustCompile(`^(?=.{1,253}\.?$)(?:(?!-|[^.]+_)[A-Za-z0-9-_]{1,63}(?<!-)(?:\.|$)){2,}$`)
	return regex.MatchString(str)
}
