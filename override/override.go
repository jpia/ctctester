package override

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OverrideRequestData struct {
	Shortcode string `json:"shortcode"`
}

func SendOverrideRequest(url, apiKey string, requestData OverrideRequestData) error {
	// JSON body
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}
