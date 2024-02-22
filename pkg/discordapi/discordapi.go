package discordapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cthulhuonice/discordchatexporter/pkg/urlbuilder"
)

/*
This is for interacting with the discord api for archiving purposes

*/

type DiscordClient struct {
	Token      string
	Bot        bool
	Guilds     []*Guild
	DMChannels []*Channel
}

type RateLimit struct {
	Message    string  `json:"message"`
	RetryAfter float32 `json:"retry_after"`
}

func (d *DiscordClient) makeRequest(uri string) (*http.Response, error) {

	// this just adds extra fmt.Println verbosity
	log_mode := false

	// create a request
	request, error := http.NewRequest(http.MethodGet, uri, nil)
	if error != nil {
		return nil, error
	}

	// fill in the token for Authorization
	token := d.Token
	if d.Bot {
		token = "Bot " + token
	}
	request.Header.Add("Authorization", token)

	// actually perform the request
	response, error := http.DefaultClient.Do(request)

	if log_mode {
		if error != nil {
			fmt.Println("makeRequest Error!")
		} else {
			fmt.Println("Request to '" + uri + "' [" + response.Status + "]")
			bodyBytes, _ := io.ReadAll(response.Body)
			fmt.Println(string(bodyBytes))
			response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	// check if we are being rate-limited
	bodyBytes, _ := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	rateLimit := RateLimit{}
	json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&rateLimit)
	if rateLimit.RetryAfter != 0.0 {
		retryAfter := rateLimit.RetryAfter
		fmt.Println("Rate limited! Retrying after", retryAfter, "seconds")
		duration := int(retryAfter * 1000)
		time.Sleep(time.Duration(duration) * time.Millisecond)
		return d.makeRequest(uri)
	}

	// move the response body back to the beginning
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return response, error
}

// use snowflake crawling to enumerate a large set of over 100 items

func (d *DiscordClient) Ping() (*http.Response, error) {
	u := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI+DISCORD_API_USER_GUILDS_URI).AddArgument("limit", "1")
	return d.makeRequest(u.BuildString())
}

func NewDiscordClient(token string, bot bool) *DiscordClient {
	d := new(DiscordClient)
	d.Token = token
	d.Bot = bot
	d.Guilds = d.EnumerateGuilds()

	response, error := d.Ping()
	if response.StatusCode != 200 {
		fmt.Println("Non-200 status code: " + response.Status)
	}
	if error != nil {
		fmt.Println("Error! ", error)
	}

	return d
}
