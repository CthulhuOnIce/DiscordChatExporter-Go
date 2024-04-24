package args

import (
	"flag"
	"fmt"
	"os"
)

const MEDIA_IGNORE = "ignore"
const MEDIA_DOWNLOAD = "download"
const MEDIA_CACHE = "cache"

// required variables
var Token string

// non-required variables that have defaults
var Dir string = "output"
var MediaDir string = "%s/media"
var BotMode bool = false
var Threads int = 1
var MediaMode string = MEDIA_IGNORE

func Init() {
	flag.StringVar(&Token, "token", "", "The token to use for the bot")
	flag.BoolVar(&BotMode, "bot", false, "Whether the token is a bot")
	flag.IntVar(&Threads, "threads", 1, "The number of threads to use")
	flag.StringVar(&MediaMode, "media", MEDIA_IGNORE, "The media mode to use (ignore, download, cache)")
	flag.StringVar(&Dir, "dir", "output", "The directory to save the output")
	flag.StringVar(&MediaDir, "mediadir", MediaDir, "The directory to save the media files in (only used if media mode is download)")
	flag.Parse()

	MediaDir = fmt.Sprintf(MediaDir, Dir)

	if MediaMode != MEDIA_IGNORE && MediaMode != MEDIA_DOWNLOAD && MediaMode != MEDIA_CACHE {
		fmt.Println("Invalid media mode")
		flag.Usage()
		os.Exit(1)
	}

	if Token == "" {
		fmt.Println("Token is required")
		flag.Usage()
		os.Exit(1)
	}

}
