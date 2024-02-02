package crawl

import (
	"fmt"
	"github.com/k3a/html2text"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CrawlAndSave(url string) (string, error) {
	// Fetch HTML content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	htmlContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Get timestamp in the specified format
	timestamp := time.Now().Format("2006-01-02-15-04")

	// Extract domain from the URL
	domain := getDomain(url)

	// Create the directory structure
	dirPath := filepath.Join("raw", domain, getPathFromURL(url), timestamp)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return dirPath, err
	}

	// Save raw.html
	htmlFilePath := filepath.Join(dirPath, "raw.html")
	if err := ioutil.WriteFile(htmlFilePath, htmlContent, os.ModePerm); err != nil {
		return dirPath, err
	}

	// Save raw.txt (rendered HTML text)
	renderedText := extractRenderedText(string(htmlContent))
	textFilePath := filepath.Join(dirPath, "raw.txt")
	if err := ioutil.WriteFile(textFilePath, []byte(renderedText), os.ModePerm); err != nil {
		return dirPath, err
	}

	fmt.Println("Files saved successfully.")
	return dirPath, nil
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
