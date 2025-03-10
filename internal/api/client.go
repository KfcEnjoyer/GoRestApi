package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	Req *Req
	Res *Res
}

func SendReq(r Req) (*http.Response, error) {
	var bodyReader io.Reader
	if r.Body != "" {
		bodyReader = strings.NewReader(r.Body)
	} else {
		bodyReader = nil
	}
	req, err := http.NewRequest(r.Method, r.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	cli := &http.Client{}
	return cli.Do(req)
}

func (c *Client) FormatResponse(resp *http.Response) *Res {
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	r := NewResponse(resp.StatusCode, resp.Header, string(read))

	return r
}
