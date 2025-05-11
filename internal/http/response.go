package http

import (
	"fmt"
	"os"
	"path/filepath"
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

func CreatedResponse() Response {
	return Response{
		HTTPVersion:  1.1,
		StatusCode:   201,
		ReasonPhrase: "Created",
		Headers:      map[string]string{},
		Body:         "",
	}
}

func InternalErrorResponse(body string) Response {
	return Response{
		HTTPVersion:  1.1,
		StatusCode:   500,
		ReasonPhrase: "Internal Server Error",
		Headers:      map[string]string{},
		Body:         body,
	}
}

func GetFileResponse(filePath string) Response {
	data, err := os.ReadFile(filepath.Join("tmp", filePath))
	if err != nil {
		return NotFoundResponse()
	}

	return OkResponse(string(data), map[string]string{})
}

func PostFileResponse(filePath string, req Request) Response {
	contentLength, exists := req.Headers["Content-Length"]
	if !exists {
		return InternalErrorResponse("Content-Length header is missing.")
	}

	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return InternalErrorResponse("Invalid Content-Length header.")
	}

	data := req.Body[:contentLengthInt]

	file, err := os.Create(filepath.Join("tmp", filePath))
	if err != nil {
		fmt.Println("Error creating file: ", err.Error())
		return InternalErrorResponse("Error creating file.")
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file: ", err.Error())
		return InternalErrorResponse("Error writing to file.")
	}

	return CreatedResponse()
}
