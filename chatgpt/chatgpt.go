package chatgpt

import (
	"encoding/json"
	"fmt"
	"github.com/jfixby/pin"
	"io/ioutil"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
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

// BuildPrompt generates a formatted prompt for ChatGPT based on rawHTML and rawText.
func BuildPrompt(rawHTML, rawText string) string {

	prompt := fmt.Sprintf(
		"Request: \n"+
			"Here is some text data that is in two sections: \n"+
			"\n"+
			"Section 1 is HTML: \n"+
			"______ \n"+
			"%s \n"+
			"______ \n"+
			"\n"+
			"Section 2 is simple text: \n"+
			"______ \n"+
			"%s \n ",
		rawHTML, rawText)

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

	var resp openai.ChatCompletionResponse
	promptChunks := splitIntoChunks(prompt, maxTokens)
	// Iterate over prompt chunks and add them as user messages

	for i, chunk := range promptChunks {
		prefix := fmt.Sprintf("Chunk[%v/%v] ", i+1, len(promptChunks))
		client.SendMessage(prefix + chunk)
	}

	// Extract the result from the response
	if len(resp.Choices) == 0 {
		return "", nil
	}
	result := resp.Choices[0].Message.Content

	return result, nil
}
