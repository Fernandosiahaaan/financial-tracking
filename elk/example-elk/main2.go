package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Struktur yang sesuai dengan response Elasticsearch
type ElasticsearchResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Timestamp time.Time `json:"timestamp"`
				Message   string    `json:"message"`
				Sensor1   float32   `json:"sensor1"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func main() {
	// URL untuk endpoint pencarian data di Elasticsearch
	elasticsearchURL := "http://localhost:9200/sensor-iot/_search"

	// Anda dapat menambahkan query string di sini untuk filter pencarian
	// Di sini kita hanya mengambil 10 data terbaru
	query := `{
		"size": 10,
		"query": {
			"match_all": {}
		},
		"sort": [
			{
				"timestamp": {
					"order": "desc"
				}
			}
		]
	}`

	// Mengirimkan request GET ke Elasticsearch dengan query JSON
	req, err := http.NewRequest("POST", elasticsearchURL, bytes.NewBuffer([]byte(query)))
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

	// Membaca response dari Elasticsearch
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Menampilkan response untuk debugging
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", body)

	// Meng-unmarshal response JSON ke dalam struktur Go
	var esResponse ElasticsearchResponse
	if err := json.Unmarshal(body, &esResponse); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Menampilkan data yang ditemukan
	for _, hit := range esResponse.Hits.Hits {
		fmt.Printf("Timestamp: %s\n", hit.Source.Timestamp)
		fmt.Printf("Message: %s\n", hit.Source.Message)
		fmt.Printf("Sensor1 Value: %.2f\n", hit.Source.Sensor1)
	}
}
