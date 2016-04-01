package conf

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Configuration is the type representing the configuration of the service
type Configuration struct {
	Listen           string
	Debug            bool
	CIURL            string
	UpdateInterval   time.Duration
	GithubOAuthToken string
	Services         []service
}

type service struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	RepoType string `yaml:"repo_type"`
	CIType   string `yaml:"ci_type"`
	Repo     string `yaml:"repo"`
	URL      string `yaml:"url"`
	Icon     string `yaml:"icon"`
	Own      bool   `yaml:"own"`
}

type unparsed struct {
	Listen           string `yaml:"listen"`
	Debug            bool   `yaml:"debug"`
	CIURL            string `yaml:"ci_url"`
	UpdateInterval   string `yaml:"update_interval"`
	GithubOAuthToken string `yaml:"github_oauth_token"`
	Services         []service
}

// C is the main configuration that is exported
var C Configuration

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string) error {
	var err error
	var u unparsed
	var d time.Duration
	conf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(conf, &u); err != nil {
		return err
	}
	if d, err = time.ParseDuration(u.UpdateInterval); err != nil {
		return err
	}
	C = Configuration{
		Listen:           u.Listen,
		Debug:            u.Debug,
		CIURL:            u.CIURL,
		UpdateInterval:   d,
		Services:         u.Services,
		GithubOAuthToken: u.GithubOAuthToken,
	}
	return nil
}
