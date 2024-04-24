package main

import (
	"github.com/cthulhuonice/discordchatexporter/pkg/args"
	"github.com/cthulhuonice/discordchatexporter/pkg/discordapi"
)

func main() {
	args.Init()

	// create new client
	client := discordapi.NewDiscordClient(args.Token, args.BotMode)

	resp, err := client.Ping()

	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		println("Non-200 status code from discord: ", resp.Status)
		panic("Failed to ping discord")
	}

	client.EnumerateGuilds()

}
