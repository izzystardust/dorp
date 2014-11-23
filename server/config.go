package main

import "github.com/BurntSushi/toml"

type Config struct {
	Key        string
	StatusPort uint16
	WebPort    uint16
}

func ReadConfig(filename string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(filename, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
