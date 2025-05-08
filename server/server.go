package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/DylanPina/go-http/server/request"
)

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

	req, err := request.Parse(string(buf))
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

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
