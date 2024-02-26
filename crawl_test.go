package main

import (
	"bufio"
	"github.com/jfixby/bchcrawler/crawl"
	"github.com/jfixby/bchcrawler/util"
	"github.com/jfixby/pin"
	"os"
	"testing"
)

func TestCrawlAndSave(t *testing.T) {
	//url := "https://github.com/mellow-finance/mellow-strategy-sdk" // Replace with your target URL
	url := "https://mellow.finance"
	dir, err := crawl.CrawlAndSave(url)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	pin.D("saved to", dir)

	data := util.ReadFile(dir + "/raw.html")
	//urls := crawl.FilterImageURLs(crawl.FindURLs(data))
	//urls := crawl.Filter(crawl.FindURLs(data),crawl.ImageFilter)
	urls := crawl.Filter(crawl.FindURLs(data), crawl.SocialNetworkFilter)
	pin.D("urls", urls)

}

func TestCrawlListAndSave(t *testing.T) {

	l, err := ReadList()
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	for _, url := range l {
		dir, err := crawl.CrawlAndSave(url)
		if err != nil {
			t.Errorf("Error occurred: %v", err)
		}

		pin.D("saved to", dir)

		data := util.ReadFile(dir + "/raw.html")
		//urls := crawl.FilterImageURLs(crawl.FindURLs(data))
		//urls := crawl.Filter(crawl.FindURLs(data),crawl.ImageFilter)
		urls := crawl.Filter(crawl.FindURLs(data), crawl.TwitterFilter)
		if len(urls) == 0 {
			continue
		}
		twitterLogoHTMLURL := urls[0] + "/photo"

		{
			dir, err := crawl.CrawlAndSave(twitterLogoHTMLURL)
			if err != nil {
				t.Errorf("Error occurred: %v", err)
			}

			pin.D("twitter saved to", dir)
		}

		pin.D("urls", urls)
	}
}

func ReadList() ([]string, error) {
	return ReadFileLines("prompts/farmap list.txt")
}

func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
