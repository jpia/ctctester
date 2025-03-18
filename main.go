package main

import (
	"math/rand"
	"time"

	"ctctester/runner"
)

func main() {
	url := "http://localhost:8080/shorten"
	overrideURL := "http://localhost:8080/admin/override/"
	userKey := "your_user_key"
	adminKey := "your_admin_key"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Number of goroutines to run
	numRunners := 5

	for i := 0; i < numRunners; i++ {
		time.Sleep(1 * time.Second)
		go runner.Run(i, url, overrideURL, userKey, adminKey, rng)
	}

	// Prevent the main function from exiting immediately
	select {}
}
