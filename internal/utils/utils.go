package utils

import (
	"os"
	"regexp"
	"strings"
)

func LoadFilterPatterns(filename string) (*regexp.Regexp, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	pattern := strings.TrimSpace(string(data))

	return regexp.MustCompile(pattern), nil
}
