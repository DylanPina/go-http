package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost/lol:4221")
	if err != nil {
		fmt.Println("Error connecting to server: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println("Error writing to server: ", err.Error())
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from server: ", err.Error())
		return
	}
	fmt.Println("Received from server: ", string(buf[:n]))
}
