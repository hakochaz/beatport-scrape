package scraper

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Holds the values for Beatport genres
var gm = map[string]string{
	"AfroHouse":                 "afro-house",
	"BassHouse":                 "bass-house",
	"BigRoom":                   "big-room",
	"Breaks":                    "breaks",
	"Dance/ElectroPop":          "dance-electro-pop",
	"DeepHouse":                 "deep-house",
	"DrumAndBass":               "drum-and-bass",
	"Dubstep":                   "dubstep",
	"ElectroHouse":              "electro-house",
	"Electronica":               "electronica",
	"Funky/Groove/Jackin'House": "funky-groove-jackin-house",
	"FutureHouse":               "future-house",
	"Garage/Bassline/Grime":     "garage-bassline-grime",
	"HardDance/Hardcore":        "hard-dance-hardcore",
	"HardTechno":                "hard-techno",
	"House":                     "house",
	"IndieDance":                "indie-dance",
	"LeftfieldBass":             "leftfield-bass",
	"LeftfieldHouseAndTechno":   "leftfield-house-and-techno",
	"MelodicHouseAndTechno":     "melodic-house-and-techno",
	"MinimalDeeptech":           "minimal-deep-tech",
	"NuDisco/Disco":             "nu-disco-disco",
	"OrganicHouseDownTempo":     "organic-house-downtempo",
	"ProgressiveHouse":          "progressive-house",
	"Psytrance":                 "psy-trance",
	"Reggae/Dancehall/Dub":      "reggae-dancehall-dub",
	"TechHouse":                 "tech-house",
	"Techno(PeakTimeDriving)":   "techno-peak-time-driving",
	"Techno(RawDeepHypnotic)":   "techno-raw-deep-hypnotic",
	"Trance":                    "trance",
	"Trap/HipHop/RAndB":         "trap-hip-hop-r-and-b",
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

// GetReleases takes in an artists list and returns a list of all tracks
// by those artists in a defined period of time
func GetReleases(al []string, co Conf) ([]Track, error) {
	if len(al) == 0 {
		return nil, errors.New("artist list is null")
	}

	g := gm[co.Genre]

	if g == "" {
		return nil, errors.New("genre not found")
	}

	var ts []Track

	c := colly.NewCollector(
		colly.Async(true),
	)

	q := make(chan Track, 5000)
	var wg sync.WaitGroup

	c.OnHTML(".horz-release-meta", func(e *colly.HTMLElement) {
		wg.Add(1)
		go func(e *colly.HTMLElement) {
			defer wg.Done()
			a := e.ChildText(".buk-horz-release-artists")
			as := strings.Split(a, ", ")
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
			e.Request.Visit("https://www.beatport.com/genre/" + g + "/1/releases?page=" + strconv.Itoa(i) + "&per-page=150&last=" + co.TimeFrame + "&type=Release")
		}
	})

	c.Visit("https://www.beatport.com/genre/" + g + "/1/releases?per-page=150&last=30d&type=Release")
	c.Wait()
	wg.Wait()
	close(q)

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
	u = "https://www.beatport.com/" + u

	return Track{a, t, l, r, u}
}
