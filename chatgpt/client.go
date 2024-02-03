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

const schema = "\"schema\": {\n    \"Name\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Name\",\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"ProfileType\": {\n      \"type\": \"options\",\n      \"name\": \"ProfileType\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"presence\": false,\n        \"inclusion\": [\n          \"Company\",\n          \"DAO\",\n          \"NFT Collection\",\n          \"Unknown\"\n        ]\n      },\n      \"optionColors\": {},\n      \"order\": 0,\n      \"visible\": true,\n      \"width\": 125\n    },\n    \"Domain\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Domain\",\n      \"order\": 4,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Short Description\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Short Description\",\n      \"order\": 5,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Long Description\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Long Description\",\n      \"order\": 6,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Twitter Link\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Twitter Link\",\n      \"order\": 7,\n      \"visible\": true,\n      \"width\": 292\n    },\n    \"Discord Link\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Discord Link\",\n      \"order\": 8,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Documentation Link\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Documentation Link\",\n      \"order\": 9,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Public Status\": {\n      \"type\": \"options\",\n      \"name\": \"Public Status\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"presence\": false,\n        \"inclusion\": [\n          \"Ded / Inactive\",\n          \"Live Outside of Bera\",\n          \"Live on TestNet\",\n          \"Officially Annouced\",\n          \"Rumoured\"\n        ]\n      },\n      \"optionColors\": {},\n      \"order\": 2,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Blog Link\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"Blog Link\",\n      \"order\": 10,\n      \"visible\": true,\n      \"width\": 200\n    },\n    \"Logo\": {\n      \"type\": \"attachment\",\n      \"name\": \"Logo\",\n      \"constraints\": {\n        \"type\": \"array\",\n        \"presence\": false\n      },\n      \"order\": 3,\n      \"visible\": true,\n      \"width\": 103\n    },\n    \"Live\": {\n      \"type\": \"boolean\",\n      \"constraints\": {\n        \"type\": \"boolean\",\n        \"presence\": false\n      },\n      \"name\": \"Live\",\n      \"order\": 11,\n      \"visible\": true,\n      \"width\": 126\n    },\n    \"onMarketMap\": {\n      \"type\": \"boolean\",\n      \"name\": \"onMarketMap\",\n      \"constraints\": {\n        \"type\": \"boolean\",\n        \"presence\": false\n      },\n      \"order\": 16,\n      \"visible\": false,\n      \"width\": 132\n    },\n    \"marketMapVersion\": {\n      \"type\": \"string\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"length\": {\n          \"maximum\": null\n        },\n        \"presence\": false\n      },\n      \"name\": \"marketMapVersion\",\n      \"order\": 17,\n      \"visible\": true,\n      \"width\": 113\n    },\n    \"MainProductType\": {\n      \"type\": \"options\",\n      \"constraints\": {\n        \"type\": \"string\",\n        \"presence\": false,\n        \"inclusion\": [\n          \"DEX\",\n          \"Game\",\n          \"Incubator\",\n          \"L1\",\n          \"Lending\",\n          \"Music\",\n          \"Name Service\",\n          \"Investor\",\n          \"Memecoin\",\n          \"Podcast\",\n          \"NFTs\",\n          \"NFT platform\",\n          \"Inscriptions\",\n          \"Aggregator\",\n          \"Launch Pad\",\n          \"Infra\"\n        ]\n      },\n      \"optionColors\": {},\n      \"name\": \"MainProductType\",\n      \"order\": 1,\n      \"visible\": true,\n      \"width\": 163\n    }\n  }"

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
