package common

import (
	"regexp"
	"strings"
)

func MatchRegex(resourceName string, filter string, regex bool) bool {
	if len(resourceName) > 0 {
		if len(filter) > 0 {
			var regexMatcher *regexp.Regexp
			if regex {
				regexMatcher = regexp.MustCompile(filter)
				return regexMatcher.MatchString(resourceName)
			} else {
				return strings.Contains(resourceName, filter)
			}

		} else if len(filter) == 0 {
			return true
		}
	}
	return false
}
