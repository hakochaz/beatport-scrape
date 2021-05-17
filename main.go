package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hakochaz/beatport-scrape/cli"
)

var ArtistsFile = "configs/artists.csv"
var OutputFile = "output/tracks.json"

func main() {
	// print welcome
	fmt.Println("Welcome to Beatport Scraper...")
	fmt.Println("Use the --help flag for more details.")

	nf := flag.Bool("new", false, "only print previously unseen releases for the timeframe")

	flag.Usage = func() {
		cli.PrintHelpMessage()
		os.Exit(0)
	}

	flag.Parse()

	fmt.Println(*nf)

	// start the cli and load the environment variables
	err := cli.StartCLI(OutputFile, ArtistsFile)

	if err != nil {
		log.Fatal("Error starting CLI...")
	}
}
