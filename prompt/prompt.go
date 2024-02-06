package prompt

import (
	"io/ioutil"
	"log"
)

type Prompt struct {
	lines []string
}

func (p *Prompt) Add(s string) {
	p.lines = append(p.lines, s)
}

func (p *Prompt) ToString() string {
	result := ""
	for _, s := range p.lines {
		result = result + s + "\n"

	}
	return result
}

func (p *Prompt) AddFile(filePath string) {
	// Read the contents of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Convert the byte slice to a string
	fileContents := string(content)
	p.Add("____")
	p.Add(fileContents)
	p.Add("____")
	p.Add("")
}

func NewPrompt() *Prompt {
	p := &Prompt{}

	return p
}
