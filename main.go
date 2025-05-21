package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func sendRequest() {
	url := "https://postman-echo.com/get?foo=bar"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading response failed: %v", err)
		return
	}

	fmt.Printf("[%s] Response:\n%s\n", time.Now().Format(time.RFC3339), body)
}

func main() {
	// Fire once immediately
	sendRequest()

	// Create a ticker that ticks every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // cleanup when main exits

	for {
		select {
		case <-ticker.C:
			sendRequest()
		}
	}
}
