package main

import (
	"log"
	"net/http"
)

const (
	listenAddr = ":8080"
)

func main() {
	http.HandleFunc("/cotacao", quoteHandler)
	log.Printf("server listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Printf("server failed: %v", err)
	}

}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	response := `{"Bid": "5.45"}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))	
}