package request

import (
	"fmt"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func (r *Request) String() string {
	return fmt.Sprintf("Method: %s\nPath: %s\nHeaders: %v\nBody: %s", r.Method, r.Path, r.Headers, r.Body)
}

func Parse(req string) (*Request, error) {
	const sep = "\r\n\r\n"

	parts := strings.SplitN(req, sep, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid HTTP request: missing header/body separator")
	}
	headerLines := strings.Split(parts[0], "\r\n")
	body := parts[1]

	reqLine := headerLines[0]
	fields := strings.SplitN(reqLine, " ", 3)
	if len(fields) < 2 {
		return nil, fmt.Errorf("Malformed request line: %q", reqLine)
	}

	headers := parseHeaders(headerLines[1:])

	return &Request{
		Method:  fields[0],
		Path:    fields[1],
		Headers: headers,
		Body:    body,
	}, nil
}

func parseHeaders(headerLines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range headerLines {
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ": ", 2)
		if len(kv) != 2 {
			continue
		}
		headers[kv[0]] = kv[1]
	}
	return headers
}
