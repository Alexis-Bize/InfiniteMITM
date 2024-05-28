package utilities

import (
	"regexp"
)

func TransformToRegex(pattern string) (*regexp.Regexp, error) {
	escapedPattern := regexp.QuoteMeta(pattern)
	escapedPattern = `^` + escapedPattern + `$`
	compiledRegex, err := regexp.Compile(escapedPattern)
	if err != nil {
		return nil, err
	}

	return compiledRegex, nil
}
