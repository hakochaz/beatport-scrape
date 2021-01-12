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
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	var as []string

	for _, a := range al {
		as = append(as, a[0])
	}

	g := os.Getenv("Genre")
	tf := os.Getenv("TimeFrame")

	tl, err := scraper.GetReleases(as, scraper.Conf{TimeFrame: tf, Genre: g})

	if err != nil {
		log.Fatal(err)
	}

	file, err := json.MarshalIndent(tl, "", " ")

	if err != nil {
		log.Fatal("Unable to marshall json")
	}

	err = ioutil.WriteFile(os.Getenv("OutputDir"), file, 0644)

	if err != nil {
		log.Fatal("Error writing tacks file")
	}
}
