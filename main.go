package main

import (
	"os"
	"strconv"

	"github.com/cthulhuonice/discordchatexporter/pkg/discordapi"
)

func main() {
	if len(os.Args) < 2 {
		panic("Need at least 1 argument")
	}

	d := discordapi.NewDiscordClient(os.Args[1], true)
	if len(os.Args) > 2 {
		guild_id, _ := strconv.Atoi(os.Args[2])
		guild, _ := d.FetchGuild(guild_id)
		println(guild.Name)

	}

}
