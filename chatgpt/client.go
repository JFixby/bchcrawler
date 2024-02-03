package chatgpt

import (
	"context"
	"errors"
	"github.com/jfixby/pin"
	openai "github.com/sashabaranov/go-openai"
	"os"
)

const model = openai.GPT3Dot5Turbo

type Client struct {
	chatgpt *openai.Client
}

func (c *Client) SendMessage(chunk string) (string, error) {
	pin.D("request >>> ", chunk)
	resp, err := c.chatgpt.CreateChatCompletion(
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
	pin.D("response <<< ", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

const initPrompt = `
You are a helpful text data crawler.
Your job is to extract information about the target blockchain project.
`

func NewClient() (*Client, error) {
	c := &Client{}

	// Retrieve OpenAI API token from environment variable
	openaiToken := os.Getenv("OPENAI_API_TOKEN")

	// Check if the token is available
	if openaiToken == "" {
		return nil, errors.New("openAI API token not set")
	}

	client := openai.NewClient(openaiToken)
	c.chatgpt = client

	pin.D("init ChatGPT request <<< ", initPrompt)

	// Create a chat completion request using the translation prompt chunks
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: initPrompt,
				},
			},
		},
	)

	pin.D("init ChatGPT response <<< ", resp.Choices[0].Message.Content)

	if err != nil {
		return nil, err
	}

	return c, nil
}
