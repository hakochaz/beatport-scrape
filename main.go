package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hakochaz/beatport-scrape/cli"
)

var ArtistsFile = os.ExpandEnv("$GOPATH/pkg/mod/github.com/hakochaz/beatport-scrape@v1.0.2/configs/artists.csv")
var OutputFile = os.ExpandEnv("$GOPATH/pkg/mod/github.com/hakochaz/beatport-scrape@v1.0.2/output/tracks.json")

func main() {
	// print welcome and help flag details
	fmt.Println("Welcome to Beatport Scraper")
	fmt.Println("Use the --help flag for more details")

	// flag for unseen
	nf := flag.Bool("new", false, "only print previously unseen releases for the timeframe")

	flag.Usage = func() {
		cli.PrintHelpMessage()
		os.Exit(0)
	}

	flag.Parse()

	args := os.Args[1:]

	var arg string
	if len(args) > 0 {
		arg = args[0]
	}

	// run the command
	if arg == "AddArtists" {
		f, err := os.OpenFile(ArtistsFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			log.Fatal("Error loading artist file")
		}

		defer f.Close()

		cli.AddArtistsPrompt(f)
	} else if arg == "SetGenre" {
		cli.SetGenrePrompt()
	} else if arg == "SetTimeframe" {
		cli.SetTimeFramePrompt()
	} else {
		// start the cli and load the environment variables
		err := cli.StartScraper(OutputFile, ArtistsFile, *nf)

		if err != nil {
			log.Fatal("Error starting CLI...")
		}
	}
}
