package discordapi

type Channel struct {
	ID       int
	Name     string
	Category string
	Guild    *Guild
	Client   *DiscordClient
	Position int
	IconURL  string   // for group DMs only
	Parent   *Channel // for threads
}

// https://discord.com/developers/docs/resources/channel

type ChannelJSON struct { // Validate this
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Type     int    `json:"type"`
	Position int    `json:"position"`
	Icon     string `json:"icon"` // group DMs only
}

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
