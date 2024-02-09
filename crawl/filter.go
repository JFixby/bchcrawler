package crawl

import (
	"net/url"
	"strings"
)

// FilterFunction defines the signature for the filtering function.
type FilterFunction func(item string) bool

// Filter applies the filtering function to each item in the list and returns a new list containing only the items that pass the filter.
func Filter(items []string, filterFn FilterFunction) []string {
	var filteredItems []string

	for _, item := range items {
		if filterFn(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func IsURLFilter(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

// Example filtering function that checks if a URL points to a social network or service.
func TwitterFilter(url string) bool {
	socialNetworks := []string{
		"twitter"}

	for _, network := range socialNetworks {
		if strings.Contains(strings.ToLower(url), network) {
			return IsURLFilter(url)
		}
	}

	return false
}

// Example filtering function that checks if a URL points to a social network or service.
func SocialNetworkFilter(url string) bool {
	socialNetworks := []string{
		"twitter",
		"facebook",
		"linkedin",
		"instagram",
		"discord",
		"reddit"}

	for _, network := range socialNetworks {
		if strings.Contains(strings.ToLower(url), network) {
			return IsURLFilter(url)
		}
	}

	return false
}

func NoFilter(url string) bool {
	return true
}

// Example filtering function that checks if a URL ends with ".png".
func ImageFilter(url string) bool {
	if strings.HasSuffix(strings.ToLower(url), ".jpg") ||
		strings.HasSuffix(strings.ToLower(url), ".jpeg") ||
		strings.HasSuffix(strings.ToLower(url), ".png") ||
		strings.HasSuffix(strings.ToLower(url), ".gif") ||
		strings.HasSuffix(strings.ToLower(url), ".bmp") ||
		strings.HasSuffix(strings.ToLower(url), ".svg") {

		return true
	}

	return false
}
