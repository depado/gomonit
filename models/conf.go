package models

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// C is the main configuration that is exported
var C Configuration

// Configuration is the type representing the configuration of the service
type Configuration struct {
	Host                   string
	Port                   string
	DefaultServiceInterval time.Duration
	DefaultRepoInterval    time.Duration
	GithubOAuthToken       string
	Debug                  bool
	Services               Services
	RawServices            []rawService
}

type rawConfiguration struct {
	Host                   string       `yaml:"host"`
	Port                   string       `yaml:"port"`
	DefaultServiceInterval string       `yaml:"default_service_interval"`
	DefaultRepoInterval    string       `yaml:"default_repo_interval"`
	GithubOAuthToken       string       `yaml:"github_oauth_token"`
	Debug                  bool         `yaml:"debug"`
	Services               []rawService `yaml:"services"`
}

func (r *rawConfiguration) ParseServices() (Services, error) {
	services := make(Services, len(r.Services))
	for i, s := range r.Services {
		p, err := s.Parse()
		if err != nil {
			return services, err
		}
		services[i] = &p
	}
	return services, nil
}

type rawService struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	URL  string `yaml:"url"`
	Icon string `yaml:"icon"`
	Own  bool   `yaml:"own"`

	RepoType  string `yaml:"repo_type"`
	RepoHost  string `yaml:"repo_url"`
	RepoToken string `yaml:"repo_token"`
	Repo      string `yaml:"repo"`

	CIType string `yaml:"ci_type"`
	CIHost string `yaml:"ci_host"`

	ServiceInterval string `yaml:"service_interval"`
	RepoInterval    string `yaml:"repo_interval"`
}

func (r *rawService) Parse() (Service, error) {
	var err error
	var s Service
	var ri time.Duration
	var si time.Duration

	if r.ID == "" || r.Name == "" {
		return s, fmt.Errorf("Configuration Error : Each service needs a 'name' and an 'id' fields.")
	}
	s.Name = r.Name
	s.ID = r.ID
	if r.RepoType != "" && r.RepoType != "github" { // TODO Modify when other types of repo are added
		return s, fmt.Errorf("Configuration Error : Service %s : %s repo type isn't supported yet", r.Name, r.RepoType)
	}
	if r.Repo != "" {
		s.Repo = r.Repo
		s.RepoType = "github" // TODO Modify when other types of repo are added
		r.RepoHost = "https://github.com"
		s.RepoURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(r.RepoHost, "/"), r.Repo)
		if r.RepoInterval != "" {
			if ri, err = time.ParseDuration(r.RepoInterval); err != nil {
				return s, fmt.Errorf("Configuration Error : Service %s : Can't parse 'repo_interval' : %s", r.Name, err)
			}
			s.RepoInterval = ri
		}
		if r.CIType != "" {
			if r.CIType != "drone" {
				return s, fmt.Errorf("Configuration Error : Service %s : %s ci type isn't supported yet", r.Name, r.CIType)
			}
			if r.CIHost == "" {
				return s, fmt.Errorf("Configuration Error : Service %s : Missing 'ci_host' field", r.Name)
			}
			s.BuildAPI = fmt.Sprintf("%s/api/repos/%s/builds", strings.TrimSuffix(r.CIHost, "/"), r.Repo)
			s.BuildURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(r.CIHost, "/"), r.Repo)
		}
	}
	if r.URL != "" {
		s.URL = r.URL
		short := strings.TrimPrefix(s.URL, "http://")
		s.ShortURL = strings.TrimPrefix(short, "https://")
		if r.ServiceInterval != "" {
			if si, err = time.ParseDuration(r.ServiceInterval); err != nil {
				return s, fmt.Errorf("Service %s : Can't parse 'service_interval' : %s", r.Name, err)
			}
			s.ServiceInterval = si
		}
	}
	s.Icon = "/static/custom/" + r.Icon
	s.Own = r.Own
	s.Host = r.Host
	return s, nil
}

// ParseConfiguration loads the given fp (file path) to the C global configuration variable.
func ParseConfiguration(fp string) error {
	var err error
	var raw rawConfiguration
	var conf []byte
	var sd time.Duration
	var rd time.Duration

	if conf, err = ioutil.ReadFile(fp); err != nil {
		return err
	}
	if err = yaml.Unmarshal(conf, &raw); err != nil {
		return err
	}
	if raw.DefaultServiceInterval == "" || raw.DefaultRepoInterval == "" {
		return fmt.Errorf("Configuration Error : Both 'default_service_interval' and 'default_repo_interval' are mandatory fields")
	}
	if sd, err = time.ParseDuration(raw.DefaultServiceInterval); err != nil {
		return fmt.Errorf("Configuration Error : Could not parse 'default_service_interval' (%s) : %s", raw.DefaultServiceInterval, err)
	}
	if rd, err = time.ParseDuration(raw.DefaultRepoInterval); err != nil {
		return fmt.Errorf("Configuration Error : Could not parse 'default_repo_interval' (%s) : %s", raw.DefaultRepoInterval, err)
	}
	ss, err := raw.ParseServices()
	if err != nil {
		return err
	}
	C = Configuration{
		Host: raw.Host,
		Port: raw.Port,
		DefaultServiceInterval: sd,
		DefaultRepoInterval:    rd,
		GithubOAuthToken:       raw.GithubOAuthToken,
		Debug:                  raw.Debug,
		Services:               ss,
	}
	All = ss
	return nil
}
