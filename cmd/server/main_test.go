package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/DylanPina/go-http/internal/client"
	"github.com/DylanPina/go-http/internal/utils"
)

// TestGet tests a basic GET request to the server
func TestGet(t *testing.T) {
	resp, err := client.ConnectHTTP("GET", "", nil, nil)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

// TestPost tests a basic GET request to the server on an unknown URL
func TestGetUnknownURL(t *testing.T) {
	resp, err := client.ConnectHTTP("GET", "/unknown", nil, nil)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}
}

// TestEchoEndpoint tests the echo endpoint
func TestEchoEndpoint(t *testing.T) {
	resp, err := client.ConnectHTTP("GET", "/echo/helloworld", nil, nil)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "helloworld" {
		t.Error("Expected body to be 'helloworld'")
	}
}

// TestReadHeader tests reading a header from the response and return it in the response body
func TestReadHeader(t *testing.T) {
	headers := map[string]string{"User-Agent": "TestClient"}
	resp, err := client.ConnectHTTP("GET", "/user-agent", headers, nil)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	headerValue := resp.Header.Get("User-Agent")

	if headerValue != headers["User-Agent"] {
		t.Errorf("Expected User-Agent to be '%s', got '%s'", headers["User-Agent"], headerValue)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "text/plain" {
		t.Error("Expected Content-Type to be 'text/plain'")
	}
}

// TestConcurrentConnections tests concurrent connections to the server
func TestConcurrentConnections(t *testing.T) {
	for i := range 10 {
		i := i // capture the current value of i

		t.Run(fmt.Sprintf("Conn-%d", i), func(t *testing.T) {
			t.Parallel() // marks this subtest as safe to run in parallel

			resp, err := client.ConnectHTTP("GET", "", nil, nil)
			if err != nil {
				t.Fatalf("Failed to make GET request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Errorf("Expected status code 200, got %d", resp.StatusCode)
			}

			_, _ = io.ReadAll(resp.Body) // ensure full body is read to enable keep-alive reuse
		})
	}
}

// TestGetFileEndpoint tests the /file endpoint and checks if the file is returned correctly
func TestGetFileEndpoint(t *testing.T) {
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	dir := projectRoot + "/tmp"
	testFileName := "testfile.txt"
	testFilePath := dir + "/" + testFileName
	testFileContent := "This is a test file."

	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	resp, err := client.ConnectHTTP("GET", "/files/"+testFileName, nil, nil)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedContent := "This is a test file."
	if string(body) != expectedContent {
		t.Errorf("Expected body to be '%s', got '%s'", expectedContent, string(body))
	}
}

// TestPostFileEndpoint tests the /file endpoint and checks if the file is created correctly
func TestPostFileEndpoint(t *testing.T) {
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	dir := projectRoot + "/tmp"
	testFileName := "testfileposted"
	testFilePath := dir + "/" + testFileName
	testFileContent := "This is a test file created from HTTP."

	headers := map[string]string{"Content-Type": "application/octet-stream", "Content-Length": fmt.Sprintf("%d", len(testFileContent))}
	resp, err := client.ConnectHTTP("POST", "/files/"+testFileName, headers, []byte(testFileContent))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()
	defer os.Remove(testFilePath) // Clean up the test file after the test

	if resp.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}

	// Read the file content from the server
	data, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(data) != testFileContent {
		t.Errorf("Expected file content to be '%s', got '%s'", testFileContent, string(data))
	}
}
