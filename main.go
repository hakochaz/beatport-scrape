package main

import (
	"fmt"
	"log"

	"github.com/hakochaz/beatport-scrape/cli"
)

var OutputFile = "configs/artists.csv"
var ArtistsFile = "output/tracks.json"

func main() {
	// Print welcome
	fmt.Println("Welcome to Beatport Scraper...")

	// start the cli and load the environment variables
	err := cli.StartCLI(OutputFile, ArtistsFile)

	if err != nil {
		log.Fatal("Error starting CLI...")
	}
}
