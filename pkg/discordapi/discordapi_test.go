package discordapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	// Set up test data
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a sample response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sample response"))
	}))
	defer mockServer.Close()

	// Set environment variable for log mode
	os.Setenv(ENV_ENABLE_REQUEST_LOG, "1")
	defer os.Unsetenv(ENV_ENABLE_REQUEST_LOG)

	// Create a mock DiscordClient
	client := &DiscordClient{
		Token: "test_token",
	}

	// Call the function being tested
	response, err := client.makeRequest(mockServer.URL)

	// Assert the expected values
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error reading response body: %v", err)
	}
	response.Body.Close()

	// Assert the expected response body
	expectedBody := "Sample response"
	if string(body) != expectedBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedBody, string(body))
	}
}

func TestMakeRequestRateLimit(t *testing.T) {
	rateLimit := new(RateLimit)
	waitTime := 1.5
	rateLimit.RetryAfter = float32(waitTime)
	totalRetries := 0
	timeoutRetries := 3
	expectedBody := "Sample response"
	var response *http.Response
	var body []byte
	var err error

	// Create a mock DiscordClient
	client := &DiscordClient{
		Token: "test_token",
	}

	// Set up test data
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a rate limit response
		totalRetries++
		if totalRetries < timeoutRetries {
			jsonData, _ := json.Marshal(rateLimit)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedBody))
		}
	}))
	defer mockServer.Close()

	os.Setenv(ENV_ENABLE_REQUEST_LOG, "1")
	defer os.Unsetenv(ENV_ENABLE_REQUEST_LOG)

	// Call the function being tested
	response, err = client.makeRequest(mockServer.URL)

	// Assert the expected values
	if totalRetries != timeoutRetries {
		t.Errorf("Expected %d retries, but got %d", timeoutRetries, totalRetries)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}

	// Read the response body
	body, err = io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unexpected error reading response body: %v", err)
	}
	response.Body.Close()

	// Assert the expected response body
	if string(body) != expectedBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedBody, string(body))
	}
}

func TestNewDiscordClient(t *testing.T) {
	// Set up test data
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a sample response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sample response"))
	}))
	defer mockServer.Close()

	const testToken = "test_token"
	const testBot = true

	DISCORD_API_BASE_URI = mockServer.URL + "/"

	// Call the function being tested
	d := NewDiscordClient(testToken, testBot)

	// Assert the expected values
	if d.Token != testToken {
		t.Errorf("Expected token '%s', but got '%s'", testToken, d.Token)
	}
	if d.Bot != testBot {
		t.Errorf("Expected bot value '%t', but got '%t'", testBot, d.Bot)
	}

	// Assert the EnumerateGuilds response
	if len(d.Guilds) != 0 {
		t.Errorf("Expected 0 guilds, but got %d", len(d.Guilds))
	}
}
