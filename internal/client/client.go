package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ConnectHTTP(method, path string, headers map[string]string, body []byte) (*http.Response, error) {
	url := "http://localhost:4221" + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}
