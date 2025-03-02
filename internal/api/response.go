package api

import (
	"fmt"
	"net/http"
)

type Res struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

func NewResponse(statusCode int, headers http.Header, body string) *Res {
	return &Res{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}
}

func (r *Res) PrintResponse() {
	fmt.Printf("Status: %d\n"+
		"Headers: %s"+
		"Body: %s", r.StatusCode, r.Headers, r.Body)
}
