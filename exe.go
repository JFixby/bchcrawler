package main

import (
	"fmt"
	"github.com/jfixby/bchcrawler/chatgpt"
)

func main() {
	timestampFolderPath := "raw/github.com/mellow-finance/mellow-strategy-sdk/2024-02-02-21-59"
	err := chatgpt.ExtratProjectDescriptionUsingChatGPT(timestampFolderPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
