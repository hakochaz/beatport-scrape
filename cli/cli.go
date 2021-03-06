package cli

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/hakochaz/beatport-scrape/scraper"
)

var envFile = os.ExpandEnv("$GOPATH/pkg/mod/github.com/hakochaz/beatport-scrape@v1.0.3/configs/env.json")

// StartScraper gets the default environment variables and
// scrapes Beatport using these, also prompts the user
// to set any variables that are currently unset
func StartScraper(outputFile, artistsFile string, new bool) error {
	var err error

	c := getEnvironmentVariables()

	if len(c.Genre) == 0 {
		fmt.Println()
		c.Genre, err = SetGenrePrompt()
		if err != nil {
			log.Fatal("Error setting genre: ", err)
		}
	}

	if len(c.TimeFrame) == 0 {
		fmt.Println()
		c.TimeFrame, err = SetTimeFramePrompt()
		if err != nil {
			log.Fatal("Error setting timeframe")
		}
	}

	// open the artists csv
	f, err := os.Open(artistsFile)

	if err != nil {
		log.Fatal("Error loading artist file: ", err)
	}

	defer f.Close()

	// get artists
	al, err := csv.NewReader(f).ReadAll()

	if err != nil {
		log.Fatal("Error loading artist file: ", err)
	}

	if len(al) == 0 {
		f, err := os.OpenFile(artistsFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			log.Fatal("Error loading artist file: ", err)
		}

		AddArtistsPrompt(f)
		f.Close()

		f, err = os.Open(artistsFile)

		if err != nil {
			log.Fatal("Error reading artists file")
		}

		defer f.Close()

		al, err = csv.NewReader(f).ReadAll()

		if err != nil {
			log.Fatal("Error reading artists file")
		}
	}

	// append artists to a slice
	var as []string

	for _, a := range al {
		if len(a) > 0 {
			as = append(as, a[0])
		}
	}

	// get releases using the scraper package
	tl, err := scraper.GetReleases(as, c)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Latest Releases: ")
	if len(tl) == 0 {
		fmt.Println("No new releases found.")
	} else if new {
		// print only unseen releases
		tj2, _ := ioutil.ReadFile("output/tracks.json")
		tl2 := make([]scraper.Track, 0)
		err := json.Unmarshal(tj2, &tl2)

		if err != nil {
			log.Fatal(err)
		}

		printUnseenReleases(tl, tl2)
	} else {
		printTracks(tl)
	}

	// marshall the json
	tj, err := json.MarshalIndent(tl, "", " ")

	if err != nil {
		log.Fatal("Unable to marshall json")
	}

	// write tracks json to the output file
	err = ioutil.WriteFile(outputFile, tj, 0644)

	if err != nil {
		log.Fatal("Error writing tacks file")
	}

	return nil
}

// SetGenrePrompt allows the user to select a default genre for the CLI
func SetGenrePrompt() (string, error) {
	fmt.Println("")
	fmt.Println("Select a default genre for scraping from the following options: ")
	fmt.Println("Enter X to exit without saving.")
	fmt.Println("")

	// Gather and sort the keys
	keys := make([]string, 0)
	for k := range scraper.GenreOptions {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		numA, _ := strconv.Atoi(keys[i])
		numB, _ := strconv.Atoi(keys[j])
		return numA < numB
	})

	for _, k := range keys {
		fmt.Println(k + " - " + scraper.GenreOptions[k][2])
	}

	scanner := bufio.NewScanner(os.Stdin)

	incorrect := false
	for {
		// reads user input until \n by default
		scanner.Scan()

		input := scanner.Text()

		if incorrect {
			fmt.Println("Please enter a valid input")
		}

		ge := scraper.GenreOptions[input]

		if input == "X" || input == "x" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		} else if len(ge) > 0 {
			return setEnvironmentVariable("Genre", input)
		}

		incorrect = true
	}
}

// SetTimeFramePrompt allows the user to select a default timeframe for the CLI
func SetTimeFramePrompt() (string, error) {
	fmt.Println("")
	fmt.Println("Select a timeframe for scraping from the following options: ")
	fmt.Println("Enter X to exit without saving.")
	fmt.Println("1 - 1 day")
	fmt.Println("2 - 7 days")
	fmt.Println("3 - 30 days")

	scanner := bufio.NewScanner(os.Stdin)

	incorrect := false
	for {
		// reads user input until \n by default
		scanner.Scan()

		input := scanner.Text()

		if incorrect {
			fmt.Println("Please enter a valid input")
		}

		if input == "1" {
			return setEnvironmentVariable("TimeFrame", "1d")
		} else if input == "2" {
			return setEnvironmentVariable("TimeFrame", "7d")
		} else if input == "3" {
			return setEnvironmentVariable("TimeFrame", "30d")
		} else if input == "X" || input == "x" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		}

		incorrect = true
	}
}

// getEnvironmentVariables returns the Conf type containing
// the environment variables
func getEnvironmentVariables() scraper.Conf {
	f, err := ioutil.ReadFile(envFile)

	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	c := scraper.Conf{}

	_ = json.Unmarshal([]byte(f), &c)

	return c
}

// setEmvironmentVariables sets an environment variable
func setEnvironmentVariable(envKey, envVal string) (string, error) {
	c := scraper.Conf{}

	err := readJSONFileIntoConf(envFile, &c)

	if err != nil {
		return "", err
	}

	if envKey == "TimeFrame" {
		c.TimeFrame = envVal
	} else if envKey == "Genre" {
		c.Genre = envVal
	}

	f, _ := json.MarshalIndent(c, "", " ")

	err = ioutil.WriteFile(envFile, f, 0644)

	return envVal, err
}

// readJSONFileIntoStruct reads a json file into a Conf type
func readJSONFileIntoConf(fn string, c *scraper.Conf) error {
	f, err := ioutil.ReadFile(fn)

	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	return json.Unmarshal([]byte(f), &c)
}

// PrintHelpMessage prints the help details if the flag is set
func PrintHelpMessage() {
	fmt.Println()
	fmt.Println("A tool for DJs to get new releases from their favourite artists")
	fmt.Println("This program allows the user to set a music genre, timeframe and a list of their favourite artists.")
	fmt.Println("The data will be used to scrape Beatport for any new releases in the specied timeframe.")
	fmt.Println()
	fmt.Println("Flags: ")
	fmt.Println()
	fmt.Println("   --new                  Only get unseeen releases since the program was last run")
	fmt.Println()
	fmt.Println("Commands: ")
	fmt.Println()
	fmt.Println("   AddArtists             Brings up a prompt for the user to enter artists to be saved to the artists file")
	fmt.Println("   SetGenre               Allows the user to set the default genre for scraping")
	fmt.Println("   SetTimeframe           Allows the user to set the default timeframe for scraping")
	fmt.Println()
}

// AddArtistsPrompt shows the user a prompt where they can add multiple
// artists to the configuration file
func AddArtistsPrompt(f *os.File) {
	fmt.Println("")
	fmt.Println("List your favourite artists using the enter key and press 1 to save.")
	fmt.Println("1 - Save File")
	fmt.Println("2 - Exit Without Saving")

	scanner := bufio.NewScanner(os.Stdin)
	al := make([][]string, 0)

	for {
		// reads user input until \n by default
		scanner.Scan()

		input := scanner.Text()

		if input == "1" {
			break
		} else if input == "2" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		} else if len(input) > 0 {
			al = append(al, []string{input})
		}
	}

	err := addArtists(al, f)

	if err != nil {
		log.Fatal("Error writing to artists file")
	}
}

// addArtists will add a new artist to the artists config file
func addArtists(a [][]string, f *os.File) error {
	// write artists
	w := csv.NewWriter(f)

	return w.WriteAll(a)
}

// printUnseenReleases print the difference between the most recent
// scrape and the last time the program was run
func printUnseenReleases(tl1, tl2 []scraper.Track) {
	nt := make([]scraper.Track, 0)

	for _, t := range tl1 {
		c := false
		for _, t2 := range tl2 {
			if t.URL == t2.URL {
				c = true
				break
			}
		}

		if !c {
			nt = append(nt, t)
		}
	}

	if len(nt) == 0 {
		fmt.Println("No new tracks have been found since the last scrape.")
	} else {
		fmt.Println("New releases since the last scrape: ")
		printTracks(nt)
	}
}

// printTracks takes a slice of Tracks and prints the values
func printTracks(t1 []scraper.Track) {
	fmt.Println()
	for _, t := range t1 {
		tj, err := json.MarshalIndent(t, "", " ")
		if err != nil {
			log.Fatal("Unable to marshall json")
		}
		fmt.Println(string(tj))
	}
	fmt.Println()
}

// CheckPermissions makes sure that the required config files are writable
func CheckFilePermissions(artistsFile, outputFile string) error {
	info, err := os.Stat(artistsFile)

	if err != nil {
		return err
	}

	m := info.Mode()

	if m != 0777 {
		err = os.Chmod(artistsFile, 0777)

		if err != nil {
			log.Println(err)
		}
	}

	info, err = os.Stat(outputFile)

	if err != nil {
		return err
	}

	m = info.Mode()

	if m != 0777 {
		err = os.Chmod(outputFile, 0777)

		if err != nil {
			log.Println(err)
		}
	}

	info, err = os.Stat(envFile)

	if err != nil {
		return err
	}

	m = info.Mode()

	if m != 0777 {
		err = os.Chmod(envFile, 0777)

		if err != nil {
			log.Println(err)
		}
	}

	return err
}
