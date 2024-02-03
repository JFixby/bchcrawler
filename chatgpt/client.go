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

const schema = `
{
  "Name": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "ProfileType": {
    "type": "options",
    "constraints": {
      "type": "string",
      "inclusion": ["Company", "DAO", "NFT Collection", "Unknown"]
    }
  },
  "Domain": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Short Description": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Long Description": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Twitter Link": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Discord Link": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Documentation Link": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Public Status": {
    "type": "options",
    "constraints": {
      "type": "string",
      "inclusion": ["Ded / Inactive", "Live Outside of Bera", "Live on TestNet", "Officially Announced", "Rumoured"]
    }
  },
  "Blog Link": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "Logo": {
    "type": "attachment",
    "constraints": {
      "type": "array"
    }
  },
  "Live": {
    "type": "boolean",
    "constraints": {
      "type": "boolean"
    }
  },
  "onMarketMap": {
    "type": "boolean",
    "constraints": {
      "type": "boolean"
    },
    "visible": false
  },
  "marketMapVersion": {
    "type": "string",
    "constraints": {
      "type": "string"
    }
  },
  "MainProductType": {
    "type": "options",
    "constraints": {
      "type": "string",
      "inclusion": ["DEX", "Game", "Incubator", "L1", "Lending", "Music", "Name Service", "Investor", "Memecoin", "Podcast", "NFTs", "NFT platform", "Inscriptions", "Aggregator", "Launch Pad", "Infra"]
    }
  }
}
`

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

const task = `
You are a helpful web crawler.
I will send you my requests.
All requests I send you will start from the word 'Request'.
The request will be sent in chunks due to its large size.
Each chunk will be marked like 'Chunk[3/5]' following the format [chunk_number/total_number_of_chunks].
I need you to acknowledge receiving each chunk with the response 'Chunk i out of N received'.

Your overall job is to:
- Extract information about the target blockchain project.
- Include information like project name, homepage, logo URL, short description (up to 350 characters),
  detailed description (up to 1000 characters), white paper link, social links, creation
  and closing dates, market symbol, GitHub, and any other relevant information.

When you receive all the chunks, I expect you to produce the resulting JSON.

`

const example = `
{
  "Name": "Example Project",
  "ProfileType": "Company",
  "Domain": "example.com",
  "Short Description": "A short project description.",
  "Long Description": "A more detailed project description.",
  "Twitter Link": "https://twitter.com/example",
  "Discord Link": "https://discord.com/example",
  "Documentation Link": "https://docs.example.com",
  "Public Status": "Live on TestNet",
  "Blog Link": "https://blog.example.com",
  "Logo": [
    {
      "url": "https://example.com/logo1.png",
      "alt": "Logo 1"
    },
    {
      "url": "https://example.com/logo2.png",
      "alt": "Logo 2"
    }
  ],
  "Live": true,
  "onMarketMap": false,
  "marketMapVersion": "v1.0",
  "MainProductType": "DEX"
}
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

	pre := task + "\n" +
		"Here is the schema for a JSON that you will need to use to produce results I request: \n"
	initPrompt := pre + schema
	pin.D("init ChatGPT request <<< ", pre)

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
