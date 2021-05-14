package cli

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hakochaz/beatport-scrape/scraper"
	"github.com/joho/godotenv"
)

// StartCLI handes the input command for the CLI
func StartCLI(outputFile, artistsFile string) error {
	var err error
	loadEnvVariables()

	g, tf := getEnvironmentVariables()

	if len(g) == 0 {
		g, err = setGenrePrompt()
		if err != nil {
			log.Fatal("Error setting genre")
		}
	}

	if len(tf) == 0 {
		tf, err = setTimeFramePrompt()
		if err != nil {
			log.Fatal("Error setting timeframe")
		}
	}

	// open the artists csv
	f, err := os.Open(artistsFile)

	if err != nil {
		log.Fatal("Error loading artist file")
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
	err = ioutil.WriteFile(outputFile, tj, 0644)

	if err != nil {
		log.Fatal("Error writing tacks file")
	}

	return nil
}

// loadEnvVariables loads the local varibales from .env
func loadEnvVariables() {
	// load the environment variables
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}
}

// setGenrePrompt allows the user to select a default genre for the CLI
func setGenrePrompt() (string, error) {
	fmt.Println("No default genre detected, please select one of the below: ")
	fmt.Println("1 - Drum And Bass")
	fmt.Println("2 - Deep House")
	fmt.Println("3 - Exit Program")

	scanner := bufio.NewScanner(os.Stdin)

	incorrect := false
	for {
		var err error
		// reads user input until \n by default
		scanner.Scan()

		input := scanner.Text()

		if incorrect {
			fmt.Println("Please enter a valid input")
		}

		if input == "1" {
			err = setEnvironmentVariable("Genre", "DrumAndBass")

			if err == nil {
				return "DrumAndBass", nil
			}
		} else if input == "2" {
			err = setEnvironmentVariable("Genre", "DeepHouse")

			if err == nil {
				return "DeepHouse", nil
			}
		} else if input == "3" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		}

		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}

		incorrect = true
	}
}

// setTimeFramePrompt allows the user to select a default timeframe for the CLI
func setTimeFramePrompt() (string, error) {
	fmt.Println("No default time frame detected, please select one of the below: ")
	fmt.Println("1 - 1 day")
	fmt.Println("2 - 7 days")
	fmt.Println("3 - 30 days")
	fmt.Println("4 - Exit Program")

	scanner := bufio.NewScanner(os.Stdin)

	incorrect := false
	for {
		var err error
		// reads user input until \n by default
		scanner.Scan()

		input := scanner.Text()

		if incorrect {
			fmt.Println("Please enter a valid input")
		}

		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}

		if input == "1" {
			err = setEnvironmentVariable("TimeFrame", "1d")

			if err == nil {
				return "1d", nil
			}
		} else if input == "2" {
			err = setEnvironmentVariable("TimeFrame", "7d")

			if err == nil {
				return "7d", nil
			}
		} else if input == "3" {
			err = setEnvironmentVariable("TimeFrame", "30d")

			if err == nil {
				return "30d", nil
			}
		} else if input == "4" {
			fmt.Println("Exiting program...")
			os.Exit(0)
		}

		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
		}

		incorrect = true
	}
}

// getEnvironmentVariables return the Genre/Timeframe variables
func getEnvironmentVariables() (string, string) {
	g := os.Getenv("Genre")
	tf := os.Getenv("TimeFrame")
	return g, tf
}

// setEmvironmentVariables sets an environment variable
func setEnvironmentVariable(envKey, envVal string) error {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()

	if err != nil {
		return err
	}

	myEnv[envKey] = envVal

	err = godotenv.Write(myEnv, "./.env")

	return err
}
