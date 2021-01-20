package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type config struct {
	Tables []string `yaml:"tables"`
}

func readConfig(filepath string) (config, error) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config{}, err
	}

	cnf := config{}
	err = yaml.Unmarshal(buf, &cnf)
	if err != nil {
		return config{}, err
	}
	return cnf, nil
}
