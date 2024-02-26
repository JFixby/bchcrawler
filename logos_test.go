package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestLogos(t *testing.T) {

	{
		url := createURL()
		jsonResponse, err := connectToEndpoint(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Response:")
		fmt.Println(formatJSON(jsonResponse))
	}
}

func createURL() string {
	// Specify the usernames that you want to lookup below
	// You can enter up to 100 comma-separated values.
	usernames := "usernames=elonmusk"
	userFields := "user.fields=description,created_at"
	// User fields are adjustable, options include:
	// created_at, description, entities, id, location, name,
	// pinned_tweet_id, profile_image_url, protected,
	// public_metrics, url, username, verified, and withheld
	url := fmt.Sprintf("https://api.twitter.com/2/users/by?%s&%s", usernames, userFields)
	return url
}

func connectToEndpoint(url string) (map[string]interface{}, error) {

	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		return nil, fmt.Errorf("missing TWITTER_BEARER_TOKEN")

	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set("User-Agent", "v2UserLookupGo")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var jsonResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned an error: %d %s", resp.StatusCode, jsonResponse)
	}

	return jsonResponse, nil
}

func formatJSON(data interface{}) string {
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes)
}
