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

	req := string(buf)
	lines := strings.Split(req, "\r\n")
	path := strings.Split(lines[0], " ")[1]

	var res string
	if path == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo/") {
		message := strings.TrimPrefix(path, "/echo/")
		res = "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: " + strconv.Itoa(len(message)) + "\r\n\r\n" +
			message
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

	fmt.Println("Response sent to: ", conn.RemoteAddr().String())
}
