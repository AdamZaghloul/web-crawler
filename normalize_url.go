package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return "", nil
	}

	normalizedURL := fmt.Sprintf("%s%s", url.Host, url.Path)

	return strings.Trim(normalizedURL, "/"), nil
}
