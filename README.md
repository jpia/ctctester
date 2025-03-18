# URL Shortener Tester

This application continuously sends POST requests to a URL shortener service at `http://localhost:8080/shorten` with random data. It tracks failed requests and stops when the number of failed requests reaches 10 or the total number of requests sent reaches 10 million. Additionally, it sends override requests to `http://localhost:8080/admin/override/:shortcode` for 20% of the future-dated successful requests from each batch.

## Assumptions

- The target URL for the URL shortener service is `http://localhost:8080/shorten`.
- The target URL for the override service is `http://localhost:8080/admin/override/:shortcode`.
- The API key for authentication is `your_user_key`.
- The admin key for override requests is `your_admin_key`.

## How It Works

1. The application generates random request data, including a `long_url` and a `release_date`.
2. It sends POST requests in batches of 100 per runner to the URL shortener service. It runs 5 concurrent runners by default.
3. If a request fails, the request data is stored in a map of failed requests.
4. The application prints the details of failed requests and the total number of requests sent after each batch.
5. The application waits for 2 seconds before sending the next batch.
6. The application selects 20% of the future-dated successful requests from each batch and sends override requests to the override service.
7. The loop continues for each runner until the number of failed requests reaches 10 or the total number of requests sent reaches 10 million.

## Usage

1. Ensure that the URL shortener service is running at `http://localhost:8080/shorten`.
2. Ensure that the override service is running at `http://localhost:8080/admin/override/:shortcode`.
3. Update the `userKey` and `adminKey` variables in `main.go` with your actual keys if different from `your_user_key` and `your_admin_key`.
4. Run the application `go run main.go`

## Dependencies

- Go 1.16 or later

## Changelog
2025-03-17: Added support for multiple runners.