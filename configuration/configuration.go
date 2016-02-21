package configuration

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Depado/gomonit/models"

	"gopkg.in/yaml.v2"
)

// Configuration is the type representing the configuration of the service
type Configuration struct {
	Listen         string
	Debug          bool
	CIURL          string
	UpdateInterval time.Duration
	Services       []service
}

type service struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	RepoType string `yaml:"repo_type"`
	CIType   string `yaml:"ci_type"`
	Repo     string `yaml:"repo"`
	URL      string `yaml:"url"`
	Icon     string `yaml:"icon"`
}

type unparsed struct {
	Listen         string `yaml:"listen"`
	Debug          bool   `yaml:"debug"`
	CIURL          string `yaml:"ci_url"`
	UpdateInterval string `yaml:"update_interval"`
	Services       []service
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
		Listen:         u.Listen,
		Debug:          u.Debug,
		CIURL:          u.CIURL,
		UpdateInterval: d,
		Services:       u.Services,
	}
	return nil
}

// Parse parses the configuration and returns the appropriate Services
func (c Configuration) Parse() (models.Services, error) {
	services := make(models.Services, len(c.Services))
	c.CIURL = strings.TrimSuffix(c.CIURL, "/")
	for i, s := range c.Services {
		if s.CIType != "" && s.CIType != "drone" {
			return nil, fmt.Errorf("Unable to use %s as CI, currently only drone is supported.", s.CIType)
		}
		if s.RepoType != "" && s.RepoType != "github" {
			return nil, fmt.Errorf("Unable to use %s as repository, currently only github is supported.", s.RepoType)
		}
		repoURL := ""
		if s.Repo != "" {
			repoURL = fmt.Sprintf("https://github.com/%s", s.Repo)
		}
		buildAPI := ""
		buildURL := ""
		if s.CIType != "" {
			buildAPI = fmt.Sprintf("%s/api/repos/%s/builds", c.CIURL, s.Repo)
			buildURL = fmt.Sprintf("%s/%s", c.CIURL, s.Repo)
		}
		short := strings.TrimPrefix(s.URL, "http://")
		short = strings.TrimPrefix(short, "https://")
		services[i] = &models.Service{
			Name:     s.Name,
			URL:      s.URL,
			ShortURL: short,
			Host:     s.Host,
			BuildAPI: buildAPI,
			BuildURL: buildURL,
			RepoURL:  repoURL,
			Icon:     "/static/custom/" + s.Icon,
		}
	}
	return services, nil
}
