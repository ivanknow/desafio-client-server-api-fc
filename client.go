package main

import (
	"context"
	"os"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"desafio-client-server-api-fc/entity"
)

const (
	serverURL  = "http://localhost:8080/cotacao"
	outputFile = "cotacao.txt"
	timeout    = 3000 * time.Millisecond
)

func main() {
	fmt.Println("Client")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()


	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("client timeout: %v", err)
		}
		log.Printf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response: %v", err)
		return
	}

	var quote entity.Quote
    err = json.Unmarshal(body, &quote)
    if err != nil {
        log.Printf("failed to unmarshal response: %v", err)
        return
    }

   
    fmt.Println("Bid:", quote.Bid)

	content := fmt.Sprintf("Dollar: %s", quote.Bid)
	if err := os.WriteFile(outputFile, []byte(content), 0o644); err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	log.Printf("saved quote to %s", outputFile)
}

