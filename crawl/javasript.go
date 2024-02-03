package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/jfixby/pin"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func CrawlAndSaveJavaScript(url string) (string, error) {
	// Create a new context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a buffer to capture the HTML content
	var htmlContent string

	// Navigate to the URL and capture the HTML content
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	); err != nil {
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
	if err := ioutil.WriteFile(htmlFilePath, []byte(htmlContent), os.ModePerm); err != nil {
		return dirPath, err
	}

	// Save raw.txt (rendered HTML text)
	renderedText := extractRenderedText(htmlContent)
	textFilePath := filepath.Join(dirPath, "raw.txt")
	if err := ioutil.WriteFile(textFilePath, []byte(renderedText), os.ModePerm); err != nil {
		return dirPath, err
	}

	pin.D("crawl js", "Files saved successfully.")
	return dirPath, nil
}
