package util

import "io/ioutil"

func ReadFile(filename string) string {
	// Read the entire file content
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Convert content to string
	strContent := string(content)

	return strContent
}
