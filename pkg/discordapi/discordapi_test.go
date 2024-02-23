package discordapi

import (
	"io/ioutil"
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
	body, err := ioutil.ReadAll(response.Body)
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
