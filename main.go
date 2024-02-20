package main

import (
	"os"

	"github.com/cthulhuonice/discordchatexporter/pkg/discordapi"
)

func main() {
	if len(os.Args) < 2 {
		panic("Need at least 1 argument")
	}

	discordapi.NewDiscordClient(os.Args[1], true)

}
