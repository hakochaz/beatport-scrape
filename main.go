package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Track holds the name/artist for the track and the Beatport url for purchasing
type Track struct {
	Artists string
	Title   string
	Labels  string
	Release string
	URL     string
}

func main() {
	f, err := os.Open("artists.csv")
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

	var ts []Track
	c := colly.NewCollector()

	c.OnHTML(".horz-release-meta", func(e *colly.HTMLElement) {
		ts = append(ts, createTrack(e))
	})

	c.OnHTML(".pagination-bottom-container", func(e *colly.HTMLElement) {
		ct := e.ChildTexts(".pag-number")
		pn, err := strconv.Atoi(ct[len(ct)-1])

		if err != nil {
			return
		}

		c.OnHTMLDetach(".pagination-bottom-container")

		for i := 2; i <= pn; i++ {
			c.Visit("https://www.beatport.com/genre/drum-and-bass/1/releases?page=" + strconv.Itoa(i) + "&per-page=150&last=30d&type=Release")
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.beatport.com/genre/drum-and-bass/1/releases?per-page=150&last=30d&type=Release")
	c.Wait()
	fmt.Println(filterArtists(ts, as))
}

func createTrack(e *colly.HTMLElement) Track {
	t := e.ChildText(".buk-horz-release-title")
	a := e.ChildText(".buk-horz-release-artists")
	a = strings.Join(strings.Fields(a), " ")
	l := e.ChildText(".buk-horz-release-labels")
	r := e.ChildText(".buk-horz-release-released")
	u := e.ChildAttr("a .buk-horz-release-title", "href")

	return Track{a, t, l, r, u}
}

func filterArtists(tl []Track, al []string) []Track {
	var ft []Track

	for _, t := range tl {
		as := strings.Split(t.Artists, ", ")
		for _, a := range as {
			for _, a2 := range al {
				if a == a2 {
					ft = append(ft, t)
				}
			}
		}
	}

	return ft
}
