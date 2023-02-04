package conf

import (
	"os"
	"time"

	"github.com/Depado/conftags"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// C is the main exported conf
var C Conf

// Conf is a configuration struct intended to be filled from a yaml file and/or
// sane defaults
type Conf struct {
	Server           Server `yaml:"server"`
	Logger           Logger `yaml:"logger"`
	GithubOAuthToken string `yaml:"github_oauth_token"`
	RServiceInterval string `yaml:"service_interval" default:"10m"`
	RRepoInterval    string `yaml:"repo_interval" default:"10m"`

	ServiceInterval time.Duration
	RepoInterval    time.Duration
	Services        []Service `yaml:"services"`
}

// Parse parses the tags and configures the logger associated
func (c *Conf) Parse() error {
	var err error

	if err = conftags.Parse(&c.Server); err != nil {
		return err
	}
	if err = conftags.Parse(c); err != nil {
		return err
	}
	c.Logger.Configure()

	if c.ServiceInterval, err = time.ParseDuration(c.RServiceInterval); err != nil {
		return errors.Wrapf(err, "configuration error: couldn't parse 'service_interval' (%s)", c.RServiceInterval)
	}
	if c.RepoInterval, err = time.ParseDuration(c.RRepoInterval); err != nil {
		return errors.Wrapf(err, "configuration error: couldn't parse 'repo_interval' (%s)", c.RRepoInterval)
	}

	return nil
}

// Load loads the configuration file into C
func Load(fp string) error {
	var err error
	var c []byte

	if c, err = os.ReadFile(fp); err != nil {
		return err
	}
	if err = yaml.Unmarshal(c, &C); err != nil {
		return err
	}

	return C.Parse()
}
