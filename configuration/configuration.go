package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Port  int `yaml:"port"`
	Hosts []struct {
		Name     string `yaml:"name"`
		URL      string `yaml:"url"`
		ShortURL string `yaml:"short_url"`
		Icon     string `yaml:"icon"`
	}
}

// C is the main configuration that is exported
var C configuration

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string) error {
	var err error
	conf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(conf, &C)
	return err
}
