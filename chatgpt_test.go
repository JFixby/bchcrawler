package main

import (
	"github.com/jfixby/pin"
	"testing"

	"github.com/jfixby/bchcrawler/chatgpt"
	"github.com/stretchr/testify/assert"
)

// Set your OpenAI access token in environment variables "OPENAI_API_TOKEN=123..."
func TestExtractProjectDescriptionJson(t *testing.T) {
	timestampFolderPath := "raw/github.com/mellow-finance/mellow-strategy-sdk/2024-02-06-15-04"
	err := chatgpt.ExtratProjectDescriptionUsingChatGPT(timestampFolderPath)
	assert.NoError(t, err, "Unexpected error")

}

//

func TestInitPrompt(t *testing.T) {
	pin.D("", chatgpt.BuildInitPrompt())

}
