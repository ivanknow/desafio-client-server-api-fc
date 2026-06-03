package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"desafio-client-server-api-fc/entity"
	_ "modernc.org/sqlite"
)

const (
	listenAddr = ":8080"
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbFile     = "cotacao.db"
	dbTimeout  = 10 * time.Millisecond
	apiTimeout = 200 * time.Millisecond
)

type apiResponse map[string]entity.Quote

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := createQuoteTable(); err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	http.HandleFunc("/cotacao", quoteHandler)
	log.Printf("server listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Printf("server failed: %v", err)
	}
}

func createQuoteTable() error {
	const sqlCreate = `CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL,
		codein TEXT NOT NULL,
		bid TEXT NOT NULL,
		ask TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err := db.Exec(sqlCreate)
	return err
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := fetchQuote(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch quote: %v", err), http.StatusBadGateway)
		return
	}

	if err := persistQuote(r.Context(), quote); err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, fmt.Sprintf("failed to persist quote: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

	log.Printf("success")
}

func fetchQuote(parentCtx context.Context) (*entity.Quote, error) {
	ctx, cancel := context.WithTimeout(parentCtx, apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("external API timeout: %v", err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	quote, ok := data["USDBRL"]
	if !ok {
		return nil, fmt.Errorf("unexpected api response")
	}

	return &quote, nil
}

func persistQuote(parentCtx context.Context, quote *entity.Quote) error {
	ctx, cancel := context.WithTimeout(parentCtx, dbTimeout)
	defer cancel()

	query := `INSERT INTO quotes (code, codein, bid, ask, timestamp, created_at) VALUES (?, ?, ?, ?, ?, datetime('now'))`
	_, err := db.ExecContext(ctx, query, quote.Code, quote.Codein, quote.Bid, quote.Ask, quote.Timestamp)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("database timeout: %v", err)
		}
	}
	return err
}
