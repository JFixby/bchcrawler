package main

import (
	"github.com/jfixby/bchcrawler/crawl"
	"github.com/jfixby/pin"
	"testing"
)

func TestCrawlAndSave(t *testing.T) {
	url := "https://github.com/mellow-finance/mellow-strategy-sdk" // Replace with your target URL
	//url := "https://mellow.finance"
	dir, err := crawl.CrawlAndSave(url)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	pin.D("saved to", dir)
}
