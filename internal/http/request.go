package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

func (r *Request) String() string {
	return fmt.Sprintf("Method: %s\nPath: %s\nProtocol: %s\nHeaders: %v\nBody: %s", r.Method, r.Path, r.Headers, r.Body)
}

func ReadConnection(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		fmt.Println("Error reading request line: ", err.Error())
		os.Exit(1)
	}

	requestLineArr := strings.Split(requestLine, " ")
	if len(requestLineArr) != 3 {
		fmt.Println("Invalid request line: ", requestLine)
		os.Exit(1)
	}

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil {
			fmt.Println("Error reading headers: ", err.Error())
			os.Exit(1)
		}

		if line == "\n\r" || line == "" {
			break
		}

		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		key := line[:colonIndex]
		value := line[colonIndex+1:]
		headers[key] = strings.TrimSpace(value)
	}

	var body string

	if value, ok := headers["Content-Length"]; ok {
		bufSize, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error converting Content-Length: ", err.Error())
			os.Exit(1)
		}
		buf := make([]byte, bufSize)
		reader.Read(buf)
		body = string(buf)
	}

	return &Request{
		Method:   requestLineArr[0],
		Path:     requestLineArr[1],
		Protocol: requestLineArr[2],
		Headers:  headers,
		Body:     body,
	}, nil
}
