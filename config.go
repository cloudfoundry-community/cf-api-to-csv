package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	APIAddress string `yaml:"apiAddress"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}

func parseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Could not open config at `%s': %s", path, err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Could not read from config (%s): %s", path, err)
	}

	conf := &Config{}
	//Read config
	err = yaml.Unmarshal(fileContents, &conf)
	if err != nil {
		return nil, fmt.Errorf("Could not parse config (%s) as YAML: %s", path, err)
	}

	return conf, nil
}

type cfConfig struct {
	AccessToken, RefreshToken string
}

func GrabToken() (*cfConfig, error) {
	file, err := os.Open(os.Getenv("HOME") + "/.cf/config.json")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var config cfConfig
	decoder.Decode(&config)
	return &config, nil
}
