package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// Use regular expression to match the IPv4 addresses.
// Ref: https://www.hackerrank.com/challenges/java-regex/forum
func IsIP(str string) bool {
	regex := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	return regex.MatchString(str)
}

// Use regular expression to match the domain name.
func IsDomainName(str string) bool {
	regex := regexp.MustCompile(`^[[:alnum:]][[:alnum:]\-]{0,61}[[:alnum:]]|[[:alpha:]]$`)
	return regex.MatchString(str)
}

func IpToBytes(ip string) []byte {
	temp := strings.Split(ip, ".")
	a, _ := strconv.Atoi(temp[0])
	b, _ := strconv.Atoi(temp[1])
	c, _ := strconv.Atoi(temp[2])
	d, _ := strconv.Atoi(temp[3])

	return []byte{
		byte(a), byte(b), byte(c), byte(d),
	}
}
