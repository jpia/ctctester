package shorten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type RequestData struct {
	LongURL     string `json:"long_url"`
	ReleaseDate string `json:"release_date"`
}

type ResponseData struct {
	Shortcode string `json:"shortcode"`
}

func SendRequest(url, apiKey string, requestData RequestData) (ResponseData, error) {
	// JSON body
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return ResponseData{}, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse response body
	var response ResponseData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ResponseData{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return response, nil
}

func GenerateRandomRequestData(rng *rand.Rand) RequestData {
	// Generate random release_date within 48 hours from now or 48 hours ago
	randomDuration := time.Duration(rng.Intn(96*60*60)-48*60*60) * time.Second
	randomReleaseDate := time.Now().Add(randomDuration).In(time.FixedZone("EST", -5*60*60)).Format(time.RFC3339)

	// Generate random path for long_url
	randomPath := fmt.Sprintf("https://example.com/%s", RandomString(rng, 10))

	return RequestData{
		LongURL:     randomPath,
		ReleaseDate: randomReleaseDate,
	}
}

func RandomString(rng *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}
