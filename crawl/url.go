package crawl

import (
	"fmt"
	"regexp"
	"strings"
)

func FindURLs(input string) []string {
	// Regular expression to match URLs
	urlRegex := regexp.MustCompile(`(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s` + "`" + `!()\[\]{};:'".,<>?«»“”‘’]))`)

	// Find all matches
	matches := urlRegex.FindAllString(input, -1)

	// Return the list of matches
	return matches
}

func FilterImageURLs(urls []string) []string {
	var imageUrls []string

	for _, url := range urls {
		// Check if the URL ends with an image extension
		if strings.HasSuffix(strings.ToLower(url), ".jpg") ||
			strings.HasSuffix(strings.ToLower(url), ".jpeg") ||
			strings.HasSuffix(strings.ToLower(url), ".png") ||
			strings.HasSuffix(strings.ToLower(url), ".gif") ||
			strings.HasSuffix(strings.ToLower(url), ".bmp") ||
			strings.HasSuffix(strings.ToLower(url), ".svg") {
			// Check if the URL is accessible
			//if isAccessible(url)
			{
				imageUrls = append(imageUrls, url)
			}
		}
	}

	return imageUrls
}

func example() {
	// Example input
	input := `Here are some links:
    - https://www.example.com
    - http://anotherexample.com/page.html
    - www.example.org
    - not a link
    `

	// Find URLs in the input
	urls := FindURLs(input)

	// Print the list of URLs
	fmt.Println("URLs found:")
	for _, url := range urls {
		fmt.Println(url)
	}
}
