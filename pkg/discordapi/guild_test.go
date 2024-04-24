package discordapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestNewGuildFromGuildJSON(t *testing.T) {
	// Set up test data
	guildJson := GuildJSON{
		ID:   "123",
		Name: "TestGuild",
		Icon: "test_icon",
	}

	// Set environment variable for log mode
	os.Setenv(ENV_ENABLE_GUILD_LOG, "1")
	defer os.Unsetenv(ENV_ENABLE_GUILD_LOG)

	// Create a mock DiscordClient
	client := &DiscordClient{}

	// Call the function being tested
	guild := NewGuildFromGuildJSON(client, guildJson)

	// Assert the expected values
	if guild.ID != 123 {
		t.Errorf("Expected guild ID to be 123, but got %d", guild.ID)
	}
	if guild.Name != "TestGuild" {
		t.Errorf("Expected guild name to be 'TestGuild', but got '%s'", guild.Name)
	}
	if guild.Icon_URL != "https://cdn.discordapp.com/icons/123/test_icon.png" {
		t.Errorf("Expected guild icon URL to be 'https://cdn.discordapp.com/icons/123/test_icon.png', but got '%s'", guild.Icon_URL)
	}

	// try again but this time with no icon
	guildJson.Icon = ""
	guild = NewGuildFromGuildJSON(client, guildJson)
	if guild.Icon_URL != "" {
		t.Errorf("Expected guild icon URL to be empty, but got '%s'", guild.Icon_URL)
	}
}

func TestFetchGuild(t *testing.T) {
	// Set up test data
	guildID := 123
	guildJson := GuildJSON{
		ID:   strconv.Itoa(guildID),
		Name: "TestGuild",
		Icon: "test_icon",
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/guilds/"+strconv.Itoa(guildID) { // Fixed the URL path
			// Respond with the guild JSON
			w.Header().Set("Content-Type", "application/json") // Set the response content type
			w.WriteHeader(http.StatusOK)                       // Set the response status code
			json.NewEncoder(w).Encode(guildJson)
		} else {
			// Respond with 404 Not Found
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Set environment variable for log mode
	os.Setenv(ENV_ENABLE_GUILD_LOG, "1")
	defer os.Unsetenv(ENV_ENABLE_GUILD_LOG)

	DISCORD_API_BASE_URI = mockServer.URL + "/"

	client := &DiscordClient{} // Initialize the client variable

	client.Bot = true
	client.Token = "test_token"

	// Call the function being tested
	guild, err := client.FetchGuild(guildID)

	// Assert the expected values
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if guild.ID != guildID {
		t.Errorf("Expected guild ID to be %d, but got %d", guildID, guild.ID)
	}
	if guild.Name != "TestGuild" {
		t.Errorf("Expected guild name to be 'TestGuild', but got '%s'", guild.Name)
	}
	if guild.Icon_URL != "https://cdn.discordapp.com/icons/123/test_icon.png" {
		t.Errorf("Expected guild icon URL to be 'https://cdn.discordapp.com/icons/123/test_icon.png', but got '%s'", guild.Icon_URL)
	}

	// Try fetching the guild again, it should return the existing guild
	guild, err = client.FetchGuild(guildID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if guild.ID != guildID {
		t.Errorf("Expected guild ID to be %d, but got %d", guildID, guild.ID)
	}
	if guild.Name != "TestGuild" {
		t.Errorf("Expected guild name to be 'TestGuild', but got '%s'", guild.Name)
	}
	if guild.Icon_URL != "https://cdn.discordapp.com/icons/123/test_icon.png" {
		t.Errorf("Expected guild icon URL to be 'https://cdn.discordapp.com/icons/123/test_icon.png', but got '%s'", guild.Icon_URL)
	}

	// Try fetching a non-existing guild, it should return an error
	nonExistingGuildID := 456
	guild, err = client.FetchGuild(nonExistingGuildID)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
	if guild != nil {
		t.Errorf("Expected guild to be nil, but got %+v", guild)
	}
}

func TestEnumerateGuilds(t *testing.T) {
	// Set up test data
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/@me/guilds" {
			// Respond with the guilds JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]GuildJSON{
				{
					ID:   "123",
					Name: "TestGuild1",
					Icon: "test_icon1",
				},
				{
					ID:   "456",
					Name: "TestGuild2",
					Icon: "test_icon2",
				},
			})
		} else if r.URL.Path == "/guilds/123" {
			// Respond with the guild JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(GuildJSON{
				ID:   "123",
				Name: "TestGuild1",
				Icon: "test_icon1",
			})
		} else if r.URL.Path == "/guilds/456" {
			// Respond with the guild JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(GuildJSON{
				ID:   "456",
				Name: "TestGuild2",
				Icon: "test_icon2",
			})
		} else {
			// Respond with 404 Not Found
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	defer mockServer.Close()

	DISCORD_API_BASE_URI = mockServer.URL + "/"

	client := &DiscordClient{} // Initialize the client variable
	client.Guilds = make(map[int]*Guild)

	// Call the function being tested
	client.EnumerateGuilds()

	// Assert the expected number of guilds
	if len(client.Guilds) != 2 {
		t.Errorf("Expected 2 guilds, but got %d", len(client.Guilds))
	}

	// Assert the expected guild values
	firstGuild, ok := client.Guilds[123]
	if !ok {
		t.Errorf("Expected guild with ID 123 to exist, but it does not")
	}
	if firstGuild.ID != 123 {
		t.Errorf("Expected guild ID to be 123, but got %d", client.Guilds[0].ID)
	}
	if firstGuild.Name != "TestGuild1" {
		t.Errorf("Expected guild name to be 'TestGuild1', but got '%s'", client.Guilds[0].Name)
	}
	if firstGuild.Icon_URL != "https://cdn.discordapp.com/icons/123/test_icon1.png" {
		t.Errorf("Expected guild icon URL to be 'https://cdn.discordapp.com/icons/123/test_icon1.png', but got '%s'", client.Guilds[0].Icon_URL)
	}

	secondGuild, ok := client.Guilds[456]

	if !ok {
		t.Errorf("Expected guild with ID 456 to exist, but it does not")
	}

	if secondGuild.ID != 456 {
		t.Errorf("Expected guild ID to be 456, but got %d", client.Guilds[1].ID)
	}
	if secondGuild.Name != "TestGuild2" {
		t.Errorf("Expected guild name to be 'TestGuild2', but got '%s'", client.Guilds[1].Name)
	}
	if secondGuild.Icon_URL != "https://cdn.discordapp.com/icons/456/test_icon2.png" {
		t.Errorf("Expected guild icon URL to be 'https://cdn.discordapp.com/icons/456/test_icon2.png', but got '%s'", client.Guilds[1].Icon_URL)
	}
}
