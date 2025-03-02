package api

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Prompt struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func CreatePrompt() *Prompt {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Method: ")
	scanner.Scan()
	method := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter the URL: ")
	scanner.Scan()
	url := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter the (JSON) Body (optional): ")
	scanner.Scan()
	body := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter headers (key:value, comma-separated, optional): ")
	scanner.Scan()
	headersInput := strings.TrimSpace(scanner.Text())

	headers := make(map[string]string)
	if headersInput != "" {
		for _, pair := range strings.Split(headersInput, ",") {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return &Prompt{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}
