package discordapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cthulhuonice/discordchatexporter/pkg/urlbuilder"
)

type Guild struct {
	ID       int
	Name     string
	Icon_URL string
	Client   *DiscordClient
}

// https://discord.com/developers/docs/resources/guild
type GuildJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

/*
func (g *Guild) FetchChannels() []string { // TODO: change to []Channel when channel type is implemented
	return u :=
}
*/

func NewGuild(id int, client *DiscordClient) *Guild {
	g := new(Guild)
	g.ID = id
	g.Client = client

	// TODO: enumerate channels

	return g
}

func NewGuildFromGuildJSON(client *DiscordClient, guildJson GuildJSON) *Guild {

	logMode := os.Getenv(ENV_ENABLE_GUILD_LOG) == "1"

	guild_id, _ := strconv.Atoi(guildJson.ID)

	guild := NewGuild(guild_id, client)
	guild.Name = guildJson.Name

	if guildJson.Icon != "" {
		guild.Icon_URL = "https://cdn.discordapp.com/icons/" + fmt.Sprint(guild_id) + "/" + guildJson.Icon + ".png"
	} else {
		guild.Icon_URL = ""
	}

	if logMode {
		fmt.Println("Creating Guild from JSON:", guild.Name, "["+fmt.Sprint(guild.ID)+"]")
	}

	return guild
}

func (d *DiscordClient) FetchGuild(guild_id int) (*Guild, error) {
	u := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI + DISCORD_API_FETCH_GUILD_URI + "/" + strconv.Itoa(guild_id))
	response, err := d.makeRequest(u.BuildString())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch guild")
	}

	// decode the json
	var guildJson GuildJSON
	if err = json.NewDecoder(response.Body).Decode(&guildJson); err != nil {
		return nil, fmt.Errorf("failed to decode guild json: %v", err)
	}

	// Check if guild already exists
	guildJsonID, err := strconv.Atoi(guildJson.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert guildJson.ID to int: %v", err)
	}

	for _, guild := range d.Guilds {
		if guild.ID == guildJsonID {
			return guild, nil
		}
	}

	guild := NewGuildFromGuildJSON(d, guildJson)

	return guild, nil
}

func (d *DiscordClient) EnumerateGuilds() {
	// 1. get all guilds with limit of 100
	//    1a. Used cached guild if in d.Guilds (map[int]*Guild)
	// 2. if guild num is < 100, immediately return with the guild list
	// 3. if guild num is >= 100, save the snowflake of the last item and repeat until
	// 		the snowflake is no longer different

	responseData := []GuildJSON{}
	var lastGuildID string = ""

	for {
		url := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI+DISCORD_API_USER_GUILDS_URI).AddArgument("limit", "100")
		if lastGuildID != "" {
			url.AddArgument("before", lastGuildID)
		}

		response, err := d.makeRequest(url.BuildString())
		if err != nil {
			fmt.Println("Failed to make request:", err)
			return
		}

		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			fmt.Println("Failed to decode guild json:", err)
			return
		}

		for i := range responseData {
			id, _ := strconv.Atoi(responseData[i].ID)
			if _, ok := d.Guilds[id]; !ok {
				guild := NewGuildFromGuildJSON(d, responseData[i])
				d.Guilds[id] = guild
			}

		}

	}

}
