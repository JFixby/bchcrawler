package main

import (
	"testing"

	"github.com/jfixby/bchcrawler/chatgpt"
	"github.com/stretchr/testify/assert"
)

// Set your OpenAI access token in environment variables "OPENAI_API_TOKEN=123..."
func TestExtractProjectDescriptionJson(t *testing.T) {
	timestampFolderPath := "raw/github.com/mellow-finance/mellow-strategy-sdk/2024-02-02-21-59"
	err := chatgpt.ExtratProjectDescriptionUsingChatGPT(timestampFolderPath)
	assert.NoError(t, err, "Unexpected error")
}
