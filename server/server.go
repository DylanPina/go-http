package main

import (
	"fmt"
	"net"
	"os"
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

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection accepted from: ", conn.RemoteAddr().String())

	buff := make([]byte, 1248)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		return
	}

	fmt.Println("Received: ", buff[:n])
}
