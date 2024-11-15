package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Struktur data log
type LogData struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type SensorData struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Sensor1   float32   `json:"sensor1"`
}

func main() {
	// Membuat data JSON yang ingin dikirim
	// logEntry := LogData{
	// 	Timestamp: time.Now(),
	// 	Level:     "info",
	// 	Message:   "This is a test log entry from Golang!",
	// }
	// Membuat data JSON yang ingin dikirim
	logEntry := SensorData{
		Timestamp: time.Now(),
		Message:   "This is a test log entry from Golang!",
		Sensor1:   10.00,
	}

	// Marshal data ke format JSON
	data, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Membuat request HTTP ke Elasticsearch
	elasticsearchURL := "http://localhost:9200/golang-logs/_doc" // URL Elasticsearch
	req, err := http.NewRequest("POST", elasticsearchURL, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Menambahkan header Content-Type: application/json
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Menampilkan response dari Elasticsearch
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Data berhasil dikirim ke Elasticsearch!")
	} else {
		fmt.Printf("Gagal mengirim data ke Elasticsearch, status: %s\n", resp.Status)
	}
}
