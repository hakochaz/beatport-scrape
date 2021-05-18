package scraper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Holds the values for Beatport genres and corresponding ids
var gm = map[string]string{
	"140-deep-dubstep-grime":    "95",
	"uk-garage-bassline":        "86",
	"afro-house":                "89",
	"bass-house":                "91",
	"big-room":                  "79",
	"breaks-breakbeat-uk-bass":  "9",
	"dance-electro-pop":         "39",
	"deep-house":                "12",
	"drum-bass":                 "1",
	"dubstep":                   "12",
	"electro-house":             "17",
	"electronica":               "3",
	"funky-groove-jackin-house": "81",
	"future-house":              "65",
	"garage-bassline-grime":     "",
	"hard-dance-hardcore":       "8",
	"hard-techno":               "2",
	"house":                     "5",
	"indie-dance":               "37",
	"melodic-house-and-techno":  "90",
	"minimal-deep-tech":         "14",
	"nu-disco-disco":            "50",
	"organic-house-downtempo":   "93",
	"progressive-house":         "15",
	"psy-trance":                "13",
	"tech-house":                "11",
	"techno-peak-time-driving":  "6",
	"techno-raw-deep-hypnotic":  "92",
	"trance":                    "7",
	"trap-wave":                 "33",
}

// Track holds the name/artist for the track and the Beatport url for purchasing
type Track struct {
	Artists string
	Title   string
	Labels  string
	Release string
	URL     string
}

// Conf holds the config data for scraping beatport tracks
type Conf struct {
	TimeFrame string
	Genre     string
}

// GetReleases takes in an artists list/config struct and returns a list of all tracks
// by those artists in a defined period of time
func GetReleases(al []string, co Conf) ([]Track, error) {
	if len(al) == 0 {
		return nil, errors.New("artist list is null")
	}

	g := co.Genre
	gn := gm[co.Genre]

	// return error if genre not found or timeframe incompatible
	if g == "" {
		return nil, errors.New("genre not found")
	}

	if co.TimeFrame != "30d" && co.TimeFrame != "7d" && co.TimeFrame != "1d" {
		return nil, errors.New("timeframe not supported")
	}

	fmt.Println("Scraping Beatport....")

	var ts []Track

	c := colly.NewCollector(
		colly.Async(true),
	)

	q := make(chan Track, 5000)
	var wg sync.WaitGroup

	// concurrenty get all the matched tracks
	c.OnHTML(".horz-release-meta", func(e *colly.HTMLElement) {
		wg.Add(1)
		go func(e *colly.HTMLElement) {
			defer wg.Done()
			as := e.ChildTexts(".buk-horz-release-artists a")

		Exit:
			for _, a := range al {
				for _, a2 := range as {
					if a == a2 {
						t := createTrack(e)
						q <- t
						continue Exit
					}
				}
			}
		}(e)
	})

	c.OnHTML(".pagination-bottom-container", func(e *colly.HTMLElement) {
		ct := e.ChildTexts(".pag-number")
		pn, err := strconv.Atoi(ct[len(ct)-1])
		if err != nil {
			return
		}

		c.OnHTMLDetach(".pagination-bottom-container")

		for i := 2; i <= pn; i++ {
			e.Request.Visit("https://www.beatport.com/genre/" + g + "/" + gn + "/releases?page=" + strconv.Itoa(i) + "&per-page=100&last=" + co.TimeFrame)
		}
	})

	c.Visit("https://www.beatport.com/genre/" + g + "/1/releases?per-page=100&last=" + co.TimeFrame)
	c.Wait()
	wg.Wait()
	close(q)

	// append the tacks into the slice
	for t := range q {
		ts = append(ts, t)
	}

	return ts, nil
}

func createTrack(e *colly.HTMLElement) Track {
	t := e.ChildText(".buk-horz-release-title")
	a := e.ChildText(".buk-horz-release-artists")
	a = strings.Join(strings.Fields(a), " ")
	l := e.ChildText(".buk-horz-release-labels")
	r := e.ChildText(".buk-horz-release-released")
	u := e.ChildAttr("p.buk-horz-release-title > a", "href")
	u = "https://www.beatport.com" + u

	return Track{a, t, l, r, u}
}
