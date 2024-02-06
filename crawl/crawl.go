package crawl

import (
	"github.com/k3a/html2text"
	"net/url"
	"strings"
)

func CrawlAndSave(url string) (string, error) {

	//return CrawlAndSaveHTML(url)
	return CrawlAndSaveJavaScript(url)
	//return CrawlAndSaveOmitter()
}

func getDomain(url string) string {
	// Extract domain from the URL
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[2]
}

func getPathFromURL(urls string) string {
	// Extract path from the URL
	u, err := url.Parse(urls)
	if err != nil {
		return ""
	}
	return u.Path
}

func extractRenderedText(htmlContent string) string {
	return html2text.HTML2Text(htmlContent)
}
