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
	u := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI + DISCORD_API_FETCH_GUILD_URI + "/" + fmt.Sprint(guild_id))
	response, error := d.makeRequest(u.BuildString())
	if error != nil {
		return nil, error
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch guild")
	}

	// decode the json
	guildJson := GuildJSON{}
	error = json.NewDecoder(response.Body).Decode(&guildJson)
	if error != nil {
		fmt.Println("Failed to decode guild json")
		panic(error)
	}

	guild := NewGuildFromGuildJSON(d, guildJson)

	d.Guilds = append(d.Guilds, guild)

	return guild, nil

}

func (d *DiscordClient) EnumerateGuilds() []*Guild {

	// 1. get all guilds with limit of 100
	// 2. if guild num is < 100, immediately return with the guild list
	// 3. if guild num is >= 100, save the snowflake of the last item and repeat until
	// 		the snowflake is no longer different

	guild_list := []*Guild{}
	responseData := []GuildJSON{}
	last_guild_id := ""

	// loop until we have all the guilds, signified by a response of less than 100
	for {
		url := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI+DISCORD_API_USER_GUILDS_URI).AddArgument("limit", "100")
		if last_guild_id != "" {
			url.AddArgument("before", last_guild_id)
		}

		response, error := d.makeRequest(url.BuildString())
		if error != nil {
			fmt.Println("Failed to make request:", error)
			return guild_list
		}

		error = json.NewDecoder(response.Body).Decode(&responseData)
		if error != nil {
			fmt.Println("Failed to decode guild json:", error)
			return guild_list
		}

		for i := range responseData {
			guild_list = append(guild_list, NewGuildFromGuildJSON(d, responseData[i]))
		}

		if len(responseData) < 100 {
			break
		}

		last_guild_id = fmt.Sprint(guild_list[len(guild_list)-1].ID)
	}

	return guild_list
}
