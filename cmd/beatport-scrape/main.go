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

	f, err := os.Open("../../configs/artists.csv")

	if err != nil {
		fmt.Print(err)
	}

	defer f.Close()

	// get artists
	al, err := csv.NewReader(f).ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	var as []string

	for _, a := range al {
		as = append(as, a[0])
	}
	fmt.Println(as)
	g := os.Getenv("Genre")
	tf := os.Getenv("TimeFrame")

	tl, _ := scraper.GetReleases(as, scraper.Conf{TimeFrame: tf, Genre: g})

	file, _ := json.MarshalIndent(tl, "", " ")
	_ = ioutil.WriteFile("../../output/tracks.json", file, 0644)
}
