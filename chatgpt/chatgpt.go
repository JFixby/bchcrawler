package chatgpt

import (
	"encoding/json"
	"fmt"
	"github.com/jfixby/bchcrawler/prompt"
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
	projectData, err := extractProjectData(string(htmlContent[:]) + string(textContent))
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

func extractProjectData(data string) (map[string]interface{}, error) {

	rawDataChunks := splitIntoChunks(data, maxTokens)

	for i, chunk := range rawDataChunks {
		prefix := fmt.Sprintf("Chunk[%v/%v] ", i+1, len(rawDataChunks))
		pin.D(prefix, len(chunk))

		{
			p := prompt.NewPrompt()
			p.Add("BEGIN RAW DATA")
			p.Add(chunk)
			p.Add("END OF RAW DATA")
			p.Add("Extract important blockchain project information from it.")
			p.Add("Return results as a pretty printed JSON.")
			p.Add("When json value is empty or null you can ignore it and exclude from the output result.")
			p.Add("Look for the following data:")
			p.AddFile("prompts/information needed.txt")

			call := p.ToString()

			result, err := sendToOpenAI(call)
			if err != nil {
				return nil, err
			}
			pin.D("result", result)

		}
	}
	//return resp, err
	//
	//prompt := BuildPrompt(rawHTML, rawText)
	//
	//result, err := sendToOpenAI(prompt)
	//if err != nil {
	//	return nil, err
	//}
	//
	//json, err := parseJSONString(result)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func BuildPrompt(rawHTML, rawText string) string {

	p := prompt.NewPrompt()
	p.Add("Discard previous data. Now we will process the next project")
	p.Add("Here is some raw text data:")
	p.Add("----")
	p.Add(rawHTML)
	p.Add("")
	p.Add("")
	p.Add(rawText)
	p.Add("----")

	p.Add("Read all the text data above.")

	p.Add("Find there the parameters and information of the project.")
	p.Add("Return to me raw pretty printed JSON.")
	p.Add("The json must include only parameters that you found.")
	p.Add("If there is no information for a specific parameter you can skip it and ignore for now.")

	return p.ToString()
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
const maxTokens = 4096 - 1000

func sendToOpenAI(chunk string) (string, error) {
	if client == nil {
		var err error
		client, err = NewClient()
		if err != nil {
			return "", err
		}
	}
	return client.SendMessage(chunk)
}
