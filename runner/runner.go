package runner

import (
	"fmt"
	"math/rand"
	"time"

	"ctctester/override"
	"ctctester/shorten"
)

func Run(workerID int, url, overrideURL, userKey, adminKey string, rng *rand.Rand) {
	fmt.Printf("Worker %d starting\n", workerID)
	failedRequests := make(map[int]shorten.RequestData)
	batchSize := 100
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
		fmt.Printf("Worker %d: Total requests sent: %d\n", workerID, totalRequestsSent)

		// Calculate and print the time of the next batch
		nextBatchTime := time.Now().Add(sleepDuration)
		fmt.Printf("Worker %d: Next batch will be sent at: %s\n", workerID, nextBatchTime.Format(time.RFC1123))

		// Handle override requests for 20% of future-dated successful requests
		handleOverrideRequests(overrideURL, adminKey, rng, successfulRequests)

		// Wait for the specified duration before sending the next batch
		time.Sleep(sleepDuration)
	}

	if len(failedRequests) >= 10 {
		fmt.Printf("Worker %d stopping as the number of failed requests has reached 10.\n", workerID)
	} else if totalRequestsSent >= maxRequests {
		fmt.Printf("Worker %d stopping as the total number of requests sent has reached 10 million.\n", workerID)
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
