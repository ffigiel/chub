package main

import (
	"encoding/json"
	"fmt"
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

func (c Config) Validate() error {
	if len(c.Commands) == 0 {
		return fmt.Errorf("no commands found")
	}
	for name, args := range c.Commands {
		if len(args) == 0 {
			return fmt.Errorf("no command specified for `%s`", name)
		}
	}
	return nil
}
