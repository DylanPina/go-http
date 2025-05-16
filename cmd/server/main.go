package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/DylanPina/go-http/internal/http"
	"github.com/DylanPina/go-http/internal/utils"
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
	for {
		req, err := http.ReadConnection(conn)
		if err != nil {
			fmt.Println("Closing connection:", conn.RemoteAddr(), "->", err)
			return
		}

		res := createResponse(*req)
		writeResponse(conn, res, *req)

		// if client asked us to close, we’ll return (defer will close)
		if val, ok := req.Headers["Connection"]; ok && val == "close" {
			fmt.Println("Client requested close:", conn.RemoteAddr())
			return
		}
	}
}

func createResponse(req http.Request) (res http.Response) {
	switch {
	case req.Path == "/":
		res = http.OkResponse("", map[string]string{})
	case strings.HasPrefix(req.Path, "/echo/"):
		message := strings.TrimPrefix(req.Path, "/echo/")
		res = http.OkResponse(message, map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": strconv.Itoa(len(message)),
		})
	case strings.HasPrefix(req.Path, "/user-agent"):
		res = http.OkResponse("", map[string]string{"User-Agent": req.Headers["User-Agent"]})
	case strings.HasPrefix(req.Path, "/files/"):
		filePath := strings.TrimPrefix(req.Path, "/files/")
		if req.Method == "GET" {
			res = http.GetFileResponse(filePath)
		} else if req.Method == "POST" {
			res = http.PostFileResponse(filePath, req)
		}
	default:
		res = http.NotFoundResponse()
	}

	return res
}

func writeResponse(conn net.Conn, res http.Response, req http.Request) {
	encoding, err := parseEncoding(req.Headers["Accept-Encoding"])
	if err != nil {
		fmt.Printf("Unsupported encodings: %s", req.Headers["Accept-Encoding"])
	}

	if encoding == "gzip" {
		err := applyCompression(&res)
		if err != nil {
			fmt.Println("Error applying compression: ", err.Error())
			return
		}
	}

	_, err = conn.Write([]byte(res.String()))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		return
	}

	fmt.Printf("Response sent to: %s\n", conn.RemoteAddr().String())
}

func parseEncoding(encoding string) (string, error) {
	if encoding == "" {
		return "", nil
	}

	encodings := strings.Split(encoding, ",")

	for _, enc := range encodings {
		if strings.TrimSpace(enc) == "gzip" {
			return "gzip", nil
		}
	}

	return "", fmt.Errorf("Unsupported encoding: %s", encoding)
}

func applyCompression(res *http.Response) error {
	compressedBody, err := utils.GzipCompress([]byte(res.Body))
	if err != nil {
		fmt.Println("Error compressing response body: ", err.Error())
		return err
	}

	res.Headers["Content-Encoding"] = "gzip"
	res.Headers["Content-Length"] = strconv.Itoa(len(compressedBody))
	res.Body = string(compressedBody)
	fmt.Println("Response body compressed with gzip")

	return nil
}
