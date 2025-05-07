package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection accepted from: ", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	req, err := parseRequest(string(buf))
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	fmt.Println("Request Method: ", req.Method)
	fmt.Println("Request Path: ", req.Path)
	fmt.Println("Request Headers: ", req.Headers)
	fmt.Println("Request Body: ", req.Body)

	var res string
	switch {
	case req.Path == "/":
		res = "HTTP/1.1 200 OK\r\n\r\n"
	case strings.HasPrefix(req.Path, "/echo/"):
		message := strings.TrimPrefix(req.Path, "/echo/")
		res = "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: " + strconv.Itoa(len(message)) + "\r\n\r\n" +
			message
	case strings.HasPrefix(req.Path, "/user-agent"):
		userAgent := req.Headers["User-Agent"]
		res = "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: " + strconv.Itoa(len(userAgent)) + "\r\n\r\n" +
			userAgent
	default:
		res = "HTTP/1.1 404 Not Found\r\n"
	}

	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

	fmt.Printf("Response sent to: %s\n", conn.RemoteAddr().String())
}

func parseRequest(req string) (*Request, error) {
	const sep = "\r\n\r\n"

	parts := strings.SplitN(req, sep, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid HTTP request: missing header/body separator")
	}
	headerLines := strings.Split(parts[0], "\r\n")
	body := parts[1]

	// request‐line
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
