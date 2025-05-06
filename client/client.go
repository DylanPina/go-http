package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":4221")
	if err != nil {
		fmt.Println("Error connecting to server: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	data := []byte("Hello, server!EOF")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error writing to server: ", err.Error())
		return
	}
}
