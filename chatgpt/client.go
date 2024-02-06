package chatgpt

import (
	"context"
	"errors"
	"github.com/jfixby/bchcrawler/prompt"
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

	initPrompt := BuildInitPrompt()

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

func BuildInitPrompt() string {
	p := prompt.NewPrompt()
	p.Add("You are a helpfull blockchain researcher")
	p.Add("Your job is to read a lot of text and search for data")
	//p.Add("Here is example of parameters you will be looking for:")
	//p.AddFile("prompts/project params.txt")
	//p.Add("Here are examples of outputs I need you to produce ideally:")
	//p.AddFile("prompts/example proect jsons.txt")
	p.Add("Is that clear? What is your job here?")

	return p.ToString()
}
