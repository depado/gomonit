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
	ServiceInterval  time.Duration
	RepoInterval     time.Duration
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
	ServiceInterval  string `yaml:"service_interval"`
	RepoInterval     string `yaml:"repo_interval"`
	GithubOAuthToken string `yaml:"github_oauth_token"`
	Services         []service
}

// C is the main configuration that is exported
var C Configuration

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string) error {
	var err error
	var u unparsed
	var sd time.Duration
	var rd time.Duration

	conf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(conf, &u); err != nil {
		return err
	}
	if sd, err = time.ParseDuration(u.ServiceInterval); err != nil {
		return err
	}
	if rd, err = time.ParseDuration(u.RepoInterval); err != nil {
		return err
	}
	C = Configuration{
		Listen:           u.Listen,
		Debug:            u.Debug,
		CIURL:            u.CIURL,
		ServiceInterval:  sd,
		RepoInterval:     rd,
		Services:         u.Services,
		GithubOAuthToken: u.GithubOAuthToken,
	}
	return nil
}
