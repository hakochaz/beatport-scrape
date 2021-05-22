package scraper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Holds the values for Beatport's genres options with corresponding
// genre with display/url format and identifier
var GenreOptions = map[string][]string{
	"1":  {"140-deep-dubstep-grime", "95", "140 / DEEP DUBSTEP / GRIME"},
	"2":  {"uk-garage-bassline", "86", "UK GARAGE / BASSLINE"},
	"3":  {"afro-house", "89", "AFRO HOUSE"},
	"4":  {"bass-house", "91", "BASS HOUSE"},
	"5":  {"big-room", "79", "BIG ROOM"},
	"6":  {"breaks-breakbeat-uk-bass", "9", "BREAKS / BREAKBEAT / UK BASS"},
	"7":  {"dance-electro-pop", "39", "DANCE / ELECTRO POP"},
	"8":  {"deep-house", "12", "DEEP HOUSE"},
	"9":  {"drum-bass", "1", "DRUM & BASS"},
	"10": {"dubstep", "12", "DUBSTEP"},
	"11": {"electro-house", "17", "ELECTRO HOUSE"},
	"12": {"electronica", "3", "ELECTRONICA"},
	"13": {"funky-groove-jackin-house", "81", "FUNKY / GROOVE / JACKIN' HOUSE"},
	"14": {"future-house", "65", "FUTURE HOUSE"},
	"15": {"hard-dance-hardcore", "8", "HARD DANCE / HARDCORE"},
	"16": {"hard-techno", "2", "HARD TECHNO"},
	"17": {"house", "5", "HOUSE"},
	"18": {"indie-dance", "37", "INDIE DANCE"},
	"19": {"melodic-house-and-techno", "90", "MELODIC HOUSE & TECHNO"},
	"20": {"minimal-deep-tech", "14", "MINIMAL / DEEP TECH"},
	"21": {"nu-disco-disco", "50", "NU DISCO / DISCO"},
	"22": {"organic-house-downtempo", "93", "ORGANIC HOUSE / DOWNTEMPO"},
	"23": {"progressive-house", "15", "PROGRESSIVE HOUSE"},
	"24": {"psy-trance", "13", "PSY-TRANCE"},
	"25": {"tech-house", "11", "TECH HOUSE"},
	"26": {"techno-peak-time-driving", "6", "TECHNO (PEAK TIME / DRIVING)"},
	"27": {"techno-raw-deep-hypnotic", "92", "TECHNO (RAW / DEEP / HYPNOTIC)"},
	"28": {"trance", "7", "TRANCE"},
	"29": {"trap-wave", "33", "TRAP / WAVE"},
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

	g := GenreOptions[co.Genre][0]
	gn := GenreOptions[co.Genre][1]

	// return error if genre not found or timeframe incompatible
	if gn == "" {
		return nil, errors.New("genre not found")
	}

	if co.TimeFrame != "30d" && co.TimeFrame != "7d" && co.TimeFrame != "1d" {
		return nil, errors.New("timeframe not supported")
	}

	fmt.Println("Scraping Beatport....")

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

	c.Visit("https://www.beatport.com/genre/" + g + "/" + gn + "/releases?per-page=100&last=" + co.TimeFrame)
	c.Wait()
	wg.Wait()
	close(q)

	tl := make([]Track, 0)

	// append the tacks into the slice
	for t := range q {
		tl = append(tl, t)
	}

	return tl, nil
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
