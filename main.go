package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hakochaz/beatport-scrape/scraper"
	"github.com/joho/godotenv"
)

func main() {
	// load the environment variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// open the artists csv
	f, err := os.Open(os.Getenv("ArtistDir"))

	if err != nil {
		fmt.Print(err)
	}

	defer f.Close()

	// get artists
	al, err := csv.NewReader(f).ReadAll()

	if err != nil {
		log.Fatal("Error loading artist file")
	}

	// append artists to a slice
	var as []string

	for _, a := range al {
		as = append(as, a[0])
	}

	g := os.Getenv("Genre")
	tf := os.Getenv("TimeFrame")

	// get releases using the scraper package
	tl, err := scraper.GetReleases(as, scraper.Conf{TimeFrame: tf, Genre: g})

	if err != nil {
		log.Fatal(err)
	}

	// marshall the json
	tj, err := json.MarshalIndent(tl, "", " ")

	if err != nil {
		log.Fatal("Unable to marshall json")
	}

	// write tracks json to the output file
	err = ioutil.WriteFile(os.Getenv("OutputDir"), tj, 0644)

	if err != nil {
		log.Fatal("Error writing tacks file")
	}
}