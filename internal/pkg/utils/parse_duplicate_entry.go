package utils

import "regexp"

func ParseDuplicateEntry(msg string) string {
	// Regular expression to capture the duplicate entry value
	re := regexp.MustCompile(`Duplicate entry '(.+?)'`)
	matches := re.FindStringSubmatch(msg)
	if len(matches) > 1 {
		return matches[1] // The first capturing group contains the duplicate value
	}
	return "unknown"
}
