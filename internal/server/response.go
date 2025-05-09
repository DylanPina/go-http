package server

import (
	"fmt"
	"os"
	"strconv"
)

type Response struct {
	HTTPVersion  float32
	StatusCode   int
	ReasonPhrase string
	Headers      map[string]string
	Body         string
}

func (r *Response) String() string {
	statusLine := fmt.Sprintf("HTTP/%.1f %d %s\r\n", r.HTTPVersion, r.StatusCode, r.ReasonPhrase)

	headers := ""

	if _, ok := r.Headers["Content-Length"]; !ok {
		r.Headers["Content-Length"] = fmt.Sprintf("%d", len(r.Body))
	}
	for key, value := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return statusLine + headers + "\r\n" + r.Body
}

func NotFoundResponse() Response {
	return Response{
		HTTPVersion:  1.1,
		StatusCode:   404,
		ReasonPhrase: "Not Found",
		Headers:      map[string]string{},
		Body:         "",
	}
}

func OkResponse(body string, headers map[string]string) Response {
	return Response{
		HTTPVersion:  1.1,
		StatusCode:   200,
		ReasonPhrase: "OK",
		Headers:      headers,
		Body:         body,
	}
}

func FileResponse(filePath string) Response {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return NotFoundResponse()
	}

	return Response{
		HTTPVersion:  1.1,
		StatusCode:   200,
		ReasonPhrase: "OK",
		Headers: map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": strconv.Itoa(len(data)),
		},
		Body: string(data),
	}
}
