package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Moyskladapiconfig
}

type Moyskladapiconfig struct {
	APIKEY     string
	TimeSpan   time.Duration
	RequestCap int
	Statehref  string
	TimeFormat string
	URLstart   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Cannot read config file")
	}

	apiKey := os.Getenv("MSAPI_KEY")
	if apiKey == "" {
		panic("API KEY does not exist")
	}

	tspnint, err := strconv.Atoi(os.Getenv("MSAPI_REQUESTCAPTIMESPAN"))
	if err != nil {
		panic("Timespan does not exist")
	}

	tspn := time.Duration(int64(tspnint)) * time.Second

	rqcap, err := strconv.Atoi(os.Getenv("MSAPI_REQUESTCAP"))
	if err != nil {
		panic("Requestcap does not exist")
	}

	stateref := os.Getenv("MSAPI_STATEHREF")
	if stateref == "" {
		panic("Statehref does not exist")
	}

	timeFormat := os.Getenv("MSAPI_TIMEFORMAT")
	if timeFormat == "" {
		panic("Timeformat does not exist")
	}

	urlstart := os.Getenv("MSAPI_URLSTART")
	if urlstart == "" {
		panic("URLstart does not exist")
	}

	return &Config{
		Moyskladapiconfig{
			APIKEY:     apiKey,
			TimeSpan:   tspn,
			RequestCap: rqcap,
			Statehref:  stateref,
			TimeFormat: timeFormat,
			URLstart:   urlstart,
		},
	}, nil
}
