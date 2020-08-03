package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Commands map[string][]string
}

func getConfig(path string) (Config, error) {
	var c Config
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(configBytes, &c)
	return c, err
}
