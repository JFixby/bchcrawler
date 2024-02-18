package twitter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//bearerToken, err := GetBearerToken()
//if err != nil {
//	fmt.Println("Error:", err)
//	return
//}
//fmt.Println("Bearer Token:", bearerToken)

func GetBearerToken() (string, error) {

	apiKey := os.Getenv("TWITTER_CONSUMER_KEY")
	apiSecretKey := os.Getenv("TWITTER_CONSUMER_SECRET")
	if apiKey == "" || apiSecretKey == "" {
		return "", fmt.Errorf("Twitter API key or API secret key not found in environment variables")
	}

	// Encode API key and API secret key
	auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecretKey))

	// Prepare request data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status code: %d", resp.StatusCode)
	}

	// Decode response body
	var tokenResponse struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", fmt.Errorf("Error decoding response body: %v", err)
	}

	if tokenResponse.TokenType != "bearer" {
		return "", fmt.Errorf("Unexpected token type: %s", tokenResponse.TokenType)
	}

	return tokenResponse.AccessToken, nil
}

func GetUserInfoAndPrint(screenName string) error {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		return fmt.Errorf("missing TWITTER_BEARER_TOKEN")
	}

	client := http.Client{}

	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.twitter.com/1.1/users/show.json?screen_name=%s", screenName), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Println("User Information:")
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
