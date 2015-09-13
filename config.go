package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ServerPort int
	DBUser     string
	DBPassword string
	DBAddress  string
	Database   string
}

func ParseConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error reading config file %s: %s", filename, err)
		return Config{}
	}
	defer file.Close()
	var config Config
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Could not read bytes from config file %s: %s", filename, err)
		return Config{}
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Printf("There was an error unmarshaling config (%s) into json: %s", filename, err)
		return Config{}
	}
	return config
}
