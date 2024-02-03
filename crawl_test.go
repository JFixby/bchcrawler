package main

import (
	"github.com/jfixby/bchcrawler/crawl"
	"testing"
)

func TestCrawlAndSave(t *testing.T) {
	//url := "https://github.com/mellow-finance/mellow-strategy-sdk" // Replace with your target URL
	url := "https://mellow.finance"
	_, err := crawl.CrawlAndSave(url)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
}
