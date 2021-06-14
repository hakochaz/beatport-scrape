# Beatport Scrape
Tool for DJs to get new releases from their favourite artists. 

Scrapes Beatport new releases from a list of artists and outputs track data in JSON format - including the Beatport URL.

Scraper configured via genre, timeframe and list of artists.

Uses Colly for the web scraping framework.

https://github.com/gocolly/colly

# Installation 
Install the program via the following command.
```
go get github.com/hakochaz/beatport-scrape
```

# Configuration
Running the application for the first time will prompt the user to set both the default Genre and Timeframe.

If the user has not added any artists to the artists file, the program will also prompt the user to list their favourite artists - and will then be appended to the configuration file.

# Commands 
The following commands allow the user to individually set the configuration varibles without scraping.

SetGenre     -     running this command will bring up the prompt which allows the user to pick from the Beatport defined lsit of genres.

SetTimeframe -     running this command will bring up a prompt to choose a default timeframe between 1, 7 and 30 days.

AddArtists   -     running this command will bring up a prompt where the user can manually list any artists to add to the artists config file.
