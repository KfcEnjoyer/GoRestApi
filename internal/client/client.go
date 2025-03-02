package client

import (
	"GoRestApi/internal/api"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	Req *api.Req
	Res *api.Res
}

func SendReq(r api.Req) (*http.Response, error) {
	var bodyReader io.Reader
	if r.Body != "" {
		bodyReader = strings.NewReader(r.Body)
	} else {
		bodyReader = nil
	}
	req, err := http.NewRequest(r.Method, r.URL, bodyReader)
	if err != nil {
		fmt.Println(err)
	}

	cli := &http.Client{}
	return cli.Do(req)
}

func (c *Client) FormatResponse(resp *http.Response) *api.Res {
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(read))

	r := api.NewResponse(resp.StatusCode, resp.Header, string(read))

	return r
}
