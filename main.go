package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	arguments "github.com/cthulhuonice/discordchatexporter/pkg/Arguments"
	"github.com/cthulhuonice/discordchatexporter/pkg/discordapi"
)

func help() {
	println("Usage: discordchatexporter -t <token> [options]")
	println("Options:")
	println("  -g  <guild_id..>  Fetch a specific guild or guilds")
	println("  -c  <channel_id..>  Fetch a specific channel or channels")
	println("  -mt <num_threads=1>  Number of channels to fetch concurrently")
	println("  -o <output_dir=output>  Directory for output")
	println("  -type <type=html-dark>  The type of logs to export (options: html-light, html-dark, txt, json, csv)")
	println("  -m  Download media files")
	println("  -b, --bot  Bot Mode")
}

func main() {
	if len(os.Args) < 2 {
		help()
		panic("Need at least 1 argument")
	}

	var token string
	var output_type string = "html-dark"
	var output_dir string = "output"
	media_download := false
	channel_targets := []string{}
	guild_targets := []string{}
	num_threads := 1
	bot_mode := false

	arguments := arguments.NewArguments()
	arguments.DebugDump()

	token = arguments.GetKey("t")
	if token == "" {
		help()
		panic("No token provided")
	}

	bot_mode = arguments.HasFlag("b", "bot")
	media_download = arguments.HasFlag("m", "media")
	num_threads, err := strconv.Atoi(arguments.GetKeyEither("mt", "multi-thread"))
	if err != nil {
		fmt.Println("Error parsing num_threads")
		panic(err)
	}
	output_dir = arguments.GetKeyEither("o", "output")
	output_type = arguments.GetKeyEither("type", "output-type")
	guild_targets = strings.Split(arguments.GetKeyEither("g", "guild"), ",")
	channel_targets = strings.Split(arguments.GetKeyEither("c", "channel"), ",")

	// TODO: make this not awful
	/*
		for arg := range os.Args {
			value := os.Args[arg]
			if value == "-t" {
				token = os.Args[arg+1]
				arg++
			} else if value == "-b" {
				bot_mode = !bot_mode
			} else if value == "-mt" {
				num_threads, _ = strconv.Atoi(os.Args[arg+1])
				arg++
			} else if value == "-o" {
				output_dir = os.Args[arg+1]
				arg++
			} else if value == "-type" {
				output_type = os.Args[arg+1]
				arg++
			} else if value == "-m" {
				media_download = true
			} else if (value == "-g") || (value == "-c") {
				targets := []string{}
				arg_len := len(os.Args)
				starting_pos := arg + 1
				for pos := range arg_len {
					if pos < starting_pos {
						continue
					}
					if string(os.Args[pos][0]) == "-" {
						break
					}
					targets = append(targets, os.Args[pos])
				}

				if value == "-g" {
					guild_targets = append(guild_targets, targets...)
				} else if value == "-c" {
					channel_targets = append(channel_targets, targets...)
				}

			}

		}
	*/

	println("Token: " + token)
	println("Threads: " + strconv.Itoa(num_threads))
	println("Media Download: ", media_download)
	println("Channel Targets: ", len(channel_targets))
	for target := range channel_targets {
		println(" - ", channel_targets[target])
	}
	println("Guild Targets: ", len(guild_targets))
	for target := range guild_targets {
		println(" - ", guild_targets[target])
	}
	println("Output Type: ", output_type)
	println("Output Dir: ", output_dir)
	println("Bot: ", bot_mode)

	client := discordapi.NewDiscordClient(token, bot_mode)
	println(len(client.Guilds))
	/*
		d := discordapi.NewDiscordClient(os.Args[1], true)
		if len(os.Args) > 2 {
			guild_id, _ := strconv.Atoi(os.Args[2])
			guild, _ := d.FetchGuild(guild_id)
			println(guild.Name)

		}
	*/

}
