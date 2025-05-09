package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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

	req, err := Parse(string(buf))
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	var res Response
	switch {
	case req.Path == "/":
		res = OkResponse("", map[string]string{})
	case strings.HasPrefix(req.Path, "/echo/"):
		message := strings.TrimPrefix(req.Path, "/echo/")
		res = OkResponse(message, map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": strconv.Itoa(len(message)),
		})
	case strings.HasPrefix(req.Path, "/user-agent"):
		res = OkResponse("", map[string]string{"User-Agent": req.Headers["User-Agent"]})
	case strings.HasPrefix(req.Path, "/files/"):
		filePath := strings.TrimPrefix(req.Path, "/files/")
		res = FileResponse(filePath)
	default:
		res = NotFoundResponse()
	}

	_, err = conn.Write([]byte(res.String()))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

	fmt.Printf("Response sent to: %s\n", conn.RemoteAddr().String())
}
