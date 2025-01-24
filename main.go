package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// httpbin.org is a service designed for testing HTTP requests
	// This endpoint returns a JSON file with the requester's IP and other info
	url := "https://httpbin.org/get"

	// Create a new request with our context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Create HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status: %s\n", resp.Status)
		return
	}

	// Create output file
	out, err := os.Create("response.json")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	// Copy data from response body to file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Successfully downloaded %d bytes\n", n)
}
