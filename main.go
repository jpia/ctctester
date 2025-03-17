package main

import (
	"fmt"
	"math/rand"
	"time"

	"ctctester/override"
	"ctctester/shorten"
)

func main() {
	url := "http://localhost:8080/shorten"
	overrideURL := "http://localhost:8080/admin/override/"
	userKey := "your_user_key"
	adminKey := "your_admin_key"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	failedRequests := make(map[int]shorten.RequestData)
	batchSize := 500
	totalRequestsSent := 0
	sleepDuration := 2 * time.Second
	maxRequests := 10_000_000

	for len(failedRequests) < 10 && totalRequestsSent < maxRequests {
		successfulRequests := make([]shorten.ResponseData, 0)

		// Send requests and track failures
		sendShortenRequests(url, userKey, rng, failedRequests, &successfulRequests, batchSize)
		totalRequestsSent += batchSize

		// Print failed requests
		printFailedRequests(failedRequests)

		// Print total requests sent
		fmt.Printf("Total requests sent: %d\n", totalRequestsSent)

		// Calculate and print the time of the next batch
		nextBatchTime := time.Now().Add(sleepDuration)
		fmt.Printf("Next batch will be sent at: %s\n", nextBatchTime.Format(time.RFC1123))

		// Handle override requests for 20% of future-dated successful requests
		handleOverrideRequests(overrideURL, adminKey, rng, successfulRequests)

		// Wait for the specified duration before sending the next batch
		time.Sleep(sleepDuration)
	}

	if len(failedRequests) >= 10 {
		fmt.Println("Stopping as the number of failed requests has reached 10.")
	} else if totalRequestsSent >= maxRequests {
		fmt.Println("Stopping as the total number of requests sent has reached 10 million.")
	}
}

func sendShortenRequests(url, userKey string, rng *rand.Rand, failedRequests map[int]shorten.RequestData, successfulRequests *[]shorten.ResponseData, batchSize int) {
	for i := 0; i < batchSize; i++ {
		// Generate random request data
		requestData := shorten.GenerateRandomRequestData(rng)

		// Send the request
		responseData, err := shorten.SendRequest(url, userKey, requestData)
		if err != nil {
			// Store failed request data in the map
			failedRequests[len(failedRequests)] = requestData
			fmt.Println("Error sending request:", err)
		} else {
			// Store successful request data
			*successfulRequests = append(*successfulRequests, responseData)
		}
	}
}

func handleOverrideRequests(url, adminKey string, rng *rand.Rand, successfulRequests []shorten.ResponseData) {
	// Select 20% of future-dated successful requests
	numOverrides := len(successfulRequests) / 5
	fmt.Printf("Sending override requests for %d successful requests\n", numOverrides)
	for i := 0; i < numOverrides; i++ {
		index := rng.Intn(len(successfulRequests))
		requestData := override.OverrideRequestData{
			Shortcode: successfulRequests[index].Shortcode,
		}
		err := override.SendOverrideRequest(url+requestData.Shortcode, adminKey, requestData)
		if err != nil {
			fmt.Println("Error sending override request:", err)
		}
	}
}

func printFailedRequests(failedRequests map[int]shorten.RequestData) {
	for i, request := range failedRequests {
		fmt.Printf("Failed Request %d: %+v\n", i, request)
	}
}
