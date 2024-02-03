package chatgpt

import (
	"encoding/json"
	"fmt"
	"github.com/jfixby/pin"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var client *Client = nil

func ExtratProjectDescriptionUsingChatGPT(timestampFolderPath string) error {
	// Read raw.html and raw.txt
	htmlFilePath := filepath.Join(timestampFolderPath, "raw.html")
	textFilePath := filepath.Join(timestampFolderPath, "raw.txt")

	htmlContent, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return err
	}

	textContent, err := ioutil.ReadFile(textFilePath)
	if err != nil {
		return err
	}

	// Extract relevant project data from HTML and text content
	projectData, err := extractProjectData(string(htmlContent[0:0]), string(textContent))
	if err != nil {
		return err
	}

	// Convert data to JSON
	jsonData, err := json.MarshalIndent(projectData, "", "  ")
	if err != nil {
		return err
	}

	// Print the pretty-printed JSON
	pin.D("json", string(jsonData))

	return nil
}

func parseJSONString(jsonString string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func extractProjectData(rawHTML string, rawText string) (map[string]interface{}, error) {

	prompt := BuildPrompt(rawHTML, rawText)

	result, err := sendToOpenAI(prompt)
	if err != nil {
		return nil, err
	}

	json, err := parseJSONString(result)
	if err != nil {
		return nil, err
	}

	return json, nil
}

const task = `
I will send you my requests.
The requests will be sent in chunks due to its large size.
 
Ypu need to find information like project name, homepage, logo URL,
short description (up to 350 characters), detailed description (up to 1000 characters),
white paper link, social links, creation and closing dates,
market symbol, GitHub, and any other relevant information that is important to know
about the project for investors.

As a result you need to to produce the resulting JSON.

Example output is:
{
  "ProjectName": "Example Project",
  "ProfileType": "Company",
  "HomePage": "example.com",
  "ShortDescription": "A short project description less than 350 characters.",
  "LongDescription": "A more detailed project description about 1000 characters.",
  "Twitter Link": "https://twitter.com/example",
  "Discord Link": "https://discord.com/example",
  "Documentation Link": "https://docs.example.com",
  "Github Link": "https://github.com/exampleproject",
  "Public Status": "Live on TestNet",
  "Blog Link": "https://blog.example.com",
  "LogoUrl": "https://example.com/logo1.png",
  "Live": true,
  "OnCoinMarketMap": false,
  "MainProductType": "DEX"
}
`

// BuildPrompt generates a formatted prompt for ChatGPT based on rawHTML and rawText.
func BuildPrompt(rawHTML, rawText string) string {

	prompt := fmt.Sprintf(
		"Request Begin: \n"+
			"Here is some text data that is in two sections: \n"+
			"\n"+
			"Section 1 is HTML: \n"+
			"______ \n"+
			"%s \n"+
			"______ \n"+
			"\n"+
			"Section 2 is simple text: \n"+
			"______ \n"+
			"%s \n "+
			"______\n"+
			"Request End.\n"+
			"\n"+
			"Task: %v \n",
		rawHTML, rawText, task)

	return prompt
}

// splitIntoChunks splits a string into chunks of the specified size
func splitIntoChunks(text string, chunkSize int) []string {
	var chunks []string
	words := strings.Fields(text)
	var currentChunk []string
	currentSize := 0

	for _, word := range words {
		wordSize := len(word) + 1 // Include space after the word
		if currentSize+wordSize > chunkSize {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{word}
			currentSize = wordSize
		} else {
			currentChunk = append(currentChunk, word)
			currentSize += wordSize
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}

// const maxTokens = 4096
const maxTokens = 4096 / 2

func sendToOpenAI(prompt string) (string, error) {

	if client == nil {
		var err error
		client, err = NewClient()
		if err != nil {
			return "", err
		}
	}

	pin.D("", "Making request to ChatGPT:")
	//pin.D("full prompt", prompt)
	pin.D("Request size is", len(prompt))

	// Split the user's prompt into smaller chunks to fit within the model's constraints
	var resp = ""
	var err error

	promptChunks := splitIntoChunks(prompt, maxTokens)
	// Iterate over prompt chunks and add them as user messages

	for i, chunk := range promptChunks {
		prefix := fmt.Sprintf("Chunk[%v/%v] ", i+1, len(promptChunks))
		resp, err = client.SendMessage(prefix + chunk)
		if err != nil {
			return "", err
		}
	}
	return resp, err
}
