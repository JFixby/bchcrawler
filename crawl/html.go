package crawl

import (
	"github.com/jfixby/pin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func CrawlAndSaveHTML(url string) (string, error) {
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

	pin.D("crawl html", "Files saved successfully.")
	return dirPath, nil
}
