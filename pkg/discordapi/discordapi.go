package discordapi

import (
	"fmt"
	"net/http"

	"github.com/cthulhuonice/discordchatexporter/pkg/urlbuilder"
)

/*
This is for interacting with the discord api for archiving purposes

*/

type DiscordClient struct {
	token string
	bot   bool
}

func (d *DiscordClient) makeRequest(uri string) (*http.Response, error) {

	// this just adds extra fmt.Println verbosity
	log_mode := true

	// create a request
	request, error := http.NewRequest(http.MethodGet, uri, nil)
	if error != nil {
		return nil, error
	}

	// fill in the token for Authorization
	token := d.token
	if d.bot {
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
		}
	}

	return response, error
}

func (d *DiscordClient) Ping() (*http.Response, error) {
	u := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI+DISCORD_API_GUILDS_URI).AddArgument("limit", "1")
	return d.makeRequest(u.BuildString())
}

func NewDiscordClient(token string, bot bool) *DiscordClient {
	d := new(DiscordClient)
	d.token = token
	d.bot = bot

	response, error := d.Ping()
	if response.StatusCode != 200 {
		fmt.Println("Non-200 status code: " + response.Status)
	}
	if error != nil {
		fmt.Println("Error! ", error)
	}
	return d
}
