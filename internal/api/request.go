package api

type Req struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func NewRequest(p *Prompt) *Req {
	return &Req{
		p.Method,
		p.URL,
		p.Headers,
		p.Body,
	}
}
