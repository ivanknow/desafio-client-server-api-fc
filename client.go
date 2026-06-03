package main

import (
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
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}	
	defer resp.Body.Close()
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

    // Access the Bid field
    fmt.Println("Bid:", quote.Bid)
}

