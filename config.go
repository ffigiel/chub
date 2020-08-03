package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Commands map[string][]string
}

func getConfig() (Config, error) {
	var c Config
	configBytes, err := ioutil.ReadFile(".chub")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(configBytes, &c)
	return c, err
}
