package discordapi

import (
	"encoding/json"
	"fmt"

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
	return g
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

	guild := NewGuild(guild_id, d)
	guild.Name = guildJson.Name

	if guildJson.Icon != "" {
		guild.Icon_URL = "https://cdn.discordapp.com/icons/" + fmt.Sprint(guild_id) + "/" + guildJson.Icon + ".png"
	} else {
		guild.Icon_URL = ""
	}

	d.Guilds = append(d.Guilds, guild)

	return guild, nil

}

func (d *DiscordClient) EnumerateGuilds() []*Guild {

	// 1. get all guilds with limit of 100
	// 2. if guild num is < 100, immediately return with the guild list
	// 3. if guild num is >= 100, save the snowflake of the last item and repeat until
	// 		the snowflake is no longer different

	guild_list := []*Guild{}
	guild_list_json := []*GuildJSON{}

	url = urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI + DISCORD_API_USER_GUILDS_URI).AddArgument("limit", "100")

	for 

}
