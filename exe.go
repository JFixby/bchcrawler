package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	timestampFolderPath := "raw/github.com/mellow-finance/mellow-strategy-sdk/2024-02-02-21-59"
	err := processTimestampFolder(timestampFolderPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func processTimestampFolder(timestampFolderPath string) error {
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
	fmt.Println(string(jsonData))

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

// BuildPrompt generates a formatted prompt for ChatGPT based on rawHTML and rawText.
func BuildPrompt(rawHTML, rawText string) string {
	task := `
- Extract information about the target blockchain project.
- Include: project name, homepage, logo URL, short description (up to 350 characters),
  detailed description (up to 1000 characters), white paper link, social links, creation
  and closing dates, market symbol, GitHub, and any other relevant information.
- Return the result as a pretty printed JSON with standardized field names.
`

	// Placeholder values
	example := map[string]interface{}{
		"name":             "Example Project",
		"homepage":         "https://example.com",
		"logoImageUrl":     "https://example.com/logo.png",
		"shortDescription": "A short project description.",
		"longDescription":  "A more detailed project description.",
		"whitePaperLink":   "https://example.com/whitepaper.pdf",
		"socialLinks": map[string]string{
			"twitter": "https://twitter.com/example",
			"github":  "https://github.com/example",
		},
		"creationDate":   "2022-01-01",
		"closingDate":    "2022-12-31",
		"marketSymbol":   "EXM",
		"additionalData": "Additional data goes here.",
	}

	// Convert data to JSON
	jsonData, _ := json.MarshalIndent(example, "", "  ")
	// Print the pretty-printed JSON
	exampleProjectJson := (string(jsonData))

	prompt := fmt.Sprintf(
		"Here is some text data that is in two sections : \n \n "+
			"Section 1 is HTML: \n %s \n"+
			"Section 2 is simple text: \n %s \n "+
			"Do the following with this data: \n %s \n"+
			"Example output: \n %s \n",
		rawHTML, rawText, task, exampleProjectJson)

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
	// Retrieve OpenAI API token from environment variable
	openaiToken := os.Getenv("OPENAI_API_TOKEN")

	// Check if the token is available
	if openaiToken == "" {
		return "", errors.New("openAI API token not set")
	}

	client := openai.NewClient(openaiToken)

	// Replace "Translate the following English text to Russian:" with the desired target language
	model := openai.GPT3Dot5Turbo

	fmt.Println("Making request to ChatGPT:")
	fmt.Println(prompt)
	fmt.Printf("Request size is %v\n", len(prompt))

	// Split the user's prompt into smaller chunks to fit within the model's constraints
	promptChunks := splitIntoChunks(prompt, maxTokens)

	// Create a chat completion request using the translation prompt chunks
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that extracts information from HTML and text content.",
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	// Iterate over prompt chunks and add them as user messages
	for _, chunk := range promptChunks {
		resp, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: chunk,
					},
				},
			},
		)
		if err != nil {
			return "", err
		}
	}

	// Extract the result from the response
	result := resp.Choices[0].Message.Content

	return result, nil
}
