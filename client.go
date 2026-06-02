package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	 fmt.Println("Client")
	req, err := http.NewRequest("GET", "http://localhost:8080/cotacao", nil)
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
	log.Printf("Response: %s", body)
}

