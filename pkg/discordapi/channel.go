package discordapi

type Channel struct {
	ID       int
	Name     string
	Category string
	Type     int
	Guild    *Guild
	Client   *DiscordClient
	Position int
	IconURL  string   // for group DMs only
	Parent   *Channel // for threads
}

// channel types https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	CHANNEL_TYPE_GUILD_TEXT          = 0
	CHANNEL_TYPE_DM                  = 1
	CHANNEL_TYPE_GUILD_VOICE         = 2
	CHANNEL_TYPE_GROUP_DM            = 3
	CHANNEL_TYPE_GUILD_CATEGORY      = 4
	CHANNEL_TYPE_GUILD_ANNOUNCEMENT  = 5
	CHANNEL_TYPE_ANNOUNCEMENT_THREAD = 10
	CHANNEL_TYPE_PUBLIC_THREAD       = 11
	CHANNEL_TYPE_PRIVATE_THREAD      = 12
	CHANNEL_TYPE_VOICE_STAGE         = 13
	CHANNEL_TYPE_GUILD_DIRECTORY     = 14
	CHANNEL_TYPE_GUILD_FORUM         = 15
	CHANNEL_TYPE_GUILD_MEDIA         = 16
)

type ChannelJSON struct { // Validate this
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Type     int    `json:"type"`
	Position int    `json:"position"`
	Icon     string `json:"icon"` // group DMs only
}

/*
func FetchChannel(c *DiscordClient, ChannelID int) *Channel { // universal channel builder, used by all the other fetchers
	url := urlbuilder.NewURLBuilder(DISCORD_API_BASE_URI + DISCORD_API_FETCH_CHANNEL_URI).BuildString()
	// make request
	response, err := c.makeRequest(url)

	if err != nil {
		panic(err) // FIXME: handle error
	}

	// parse response
	var channelJSON ChannelJSON
	error := json.NewDecoder(response.Body).Decode(&channelJSON)
	if error != nil {
		panic(error) // FIXME: handle error
	}

	channel := new(Channel)

	// convert channelJSON to channel
	channel.ID, _ = strconv.Atoi(channelJSON.ID)
	channel.Name = channelJSON.Name
	channel.Category = channelJSON.Category
	channel.Type = channelJSON.Type
	channel.Client = c
	channel.Position = channelJSON.Position

	if channelJSON.Icon != "" {
		channel.IconURL = "https://cdn.discordapp.com/icons/" + channelJSON.ID + "/" + channelJSON.Icon + ".png"
	} else {
		channel.IconURL = ""
	}

	return channel
}

// https://discord.com/developers/docs/resources/channel#get-channel-messages


func (c *Channel) FetchThread(ThreadID int) []*Channel {
	return
}

func (c *Channel) EnumerateThreads() []*Channel {
	return
}

func (c *DiscordClient) FetchDMChannel(ChannelID int) *Channel {

}

func (c *DiscordClient) EnumerateDMChannels() []*Channel {

}

func (d *Guild) FetchChannel(ChannelID int) *Channel {
	return
}

func (d *Guild) EnumerateChannels() []*Channel {
	return
}
*/
