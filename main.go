package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/md4"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://httpbin.org/get"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status: %s\n", resp.Status)
		return
	}

	// Create an MD4 hasher
	hasher := md4.New()

	// Create output file
	out, err := os.Create("response.json")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	// Use MultiWriter to write to both the file and the hasher
	writer := io.MultiWriter(out, hasher)

	// Copy data from response body to both writers
	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	// Get the hash sum
	hashSum := hasher.Sum(nil)

	fmt.Printf("Successfully downloaded %d bytes\n", n)
	fmt.Printf("MD4 hash of the content: %x\n", hashSum)
}
