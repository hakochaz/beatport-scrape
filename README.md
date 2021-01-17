# Beatport Scrape
Tool for DJs to get new releases from their favourite artists. 

Scrapes Beatport releases from a list of artists and outputs track data - including the Beatport URL.

Scraper configured via genre and timeframe.

Uses Colly for the web scraping framework.

https://github.com/gocolly/colly

# Configuration 
Start by adding your favourite artists to the artists csv file located in the configs folder.

Edit the .Env file in the root directory.

1. Timeframe

Choose any of the following timeframes: 
1d, 7d or 30d.

Default: 30d

2. Genre

Supports all availble Beatport genres, choose from the following:

AfroHouse, BassHouse, BigRoom, Breaks, Dance/ElectroPop, DeepHouse, DrumAndBass, Dubstep, ElectroHouse,
Electronica, Funky/Groove/Jackin'House, FutureHouse Garage/Bassline/Grime, HardDance/Hardcore, HardTechno,
House, IndieDance, LeftfieldBass, LeftfieldHouseAndTechno, MelodicHouseAndTechno, MinimalDeeptech,
NuDisco/Disco, OrganicHouseDownTempo, ProgressiveHouse, Psytrance, Reggae/Dancehall/Dub, TechHouse,
Techno(PeakTimeDriving), Techno(RawDeepHypnotic), Trance, Trap/HipHop/RAndB.

Default: DrumAndBass

3. ArtistDir

The directory which holds the artists csv file.

Default: configs/artists.csv

4. OutputDir

The directory to output any matching tracks in JSON format.

Default: output/tracks.json
