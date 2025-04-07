package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/depado/gomonit/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// All represents all the services
var All Services

// Repo holds information on a repository
type Repo struct {
	URL  string `json:"url"`
	Path string `json:"path"`
	Type string `json:"type"`
	Host string `json:"host"`

	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Watchers    int    `json:"watchers"`
}

// CI represents data from the CI
type CI struct {
	API string `json:"api"`
	URL string `json:"url"`
}

// Service is a single service
type Service struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	ShortURL        string        `json:"short_url"`
	Host            string        `json:"host"`
	ServiceInterval time.Duration `json:"service_interval"`

	Repo *Repo `json:"repo,omitempty"`
	CI   *CI   `json:"ci,omitempty"`

	Last            string        `json:"last"`
	RespTime        time.Duration `json:"resp_time"`
	Status          int           `json:"status"`
	Icon            string        `json:"icon"`
	CurrentBuildURL string        `json:"current_build"`
	LastBuilds      Builds        `json:"last_builds"`
	LastCommits     Commits       `json:"last_commits"`
	Own             bool          `json:"own"`
}

// InitializeServices grabs all the services from the configuration and
// initializes the All variable
func InitializeServices() error {
	var err error
	All, err = ParseServicesFromConf(conf.C)
	return err
}

// ParseServicesFromConf parses all the services in the configuration struct
// and returns a slice of pointers to Service
func ParseServicesFromConf(c conf.Conf) (Services, error) {
	var err error
	var s *Service
	var ss Services

	for _, cs := range c.Services {
		if s, err = NewServiceFromConf(cs); err != nil {
			return ss, errors.Wrap(err, "parse all services")
		}
		ss = append(ss, s)
	}

	return ss, nil
}

// NewServiceFromConf parses a configured service and returns a service
func NewServiceFromConf(cs conf.Service) (*Service, error) {
	s := Service{
		Name: cs.Name,
		Icon: "/static/custom/" + cs.Icon,
		Own:  cs.Own,
		Host: cs.Host,
	}

	if s.Name == "" {
		return &s, fmt.Errorf("configuration error: each service needs a 'name' field")
	}
	if cs.Repo != nil {
		if cs.Repo.Type != "github" {
			return &s, fmt.Errorf("configuration error: service %s - %s repo type isn't supported", cs.Name, cs.Repo.Type)
		}
		s.Repo = &Repo{
			Path: cs.Repo.Path,
			Type: cs.Repo.Type,
		}
		switch s.Repo.Type {
		case "github":
			s.Repo.Host = "https://github.com"
			s.Repo.URL = fmt.Sprintf("%s/%s", strings.TrimSuffix(cs.Repo.Host, "/"), cs.Repo.Path)
		}
		if cs.CI != nil {
			if cs.CI.Type != "drone" {
				return &s, fmt.Errorf("configuration error: service %s - %s ci type isn't supported", cs.Name, cs.CI.Type)
			}
			if cs.CI.Host == "" {
				return &s, fmt.Errorf("configuration error: service %s - ci missing 'host' field", cs.Name)
			}
			s.CI = &CI{
				API: fmt.Sprintf("%s/api/repos/%s/builds", strings.TrimSuffix(cs.CI.Host, "/"), cs.Repo.Path),
				URL: fmt.Sprintf("%s/%s", strings.TrimSuffix(cs.CI.Host, "/"), cs.Repo.Path),
			}
		}
	}
	if cs.URL != "" {
		s.URL = cs.URL
		short := strings.TrimPrefix(cs.URL, "http://")
		s.ShortURL = strings.TrimPrefix(short, "https://")
	}

	return &s, nil
}

// FetchStatus checks if the service is running
func (s *Service) FetchStatus() {
	clog := logrus.WithFields(logrus.Fields{"action": "status", "service": s.Name})
	tp := newTransport()
	client := &http.Client{Transport: tp}
	start := time.Now()
	s.Last = start.Format("2006/01/02 15:04:05")

	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		clog.WithError(err).Error("Couldn't create request")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		clog.WithError(err).Warn("Couldn't fetch status")
		return
	}
	defer resp.Body.Close()        //nolint:errcheck
	io.Copy(io.Discard, resp.Body) //nolint:errcheck

	d := tp.ReqDuration()
	s.RespTime = d - (d % time.Millisecond)
	s.Status = resp.StatusCode
}

// FetchBuilds checks the last build
func (s *Service) FetchBuilds() {
	resp, err := http.Get(s.CI.API)
	if err != nil {
		log.Printf("[%s][ERROR] While requesting build status : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close() //nolint:errcheck
	var all UnparsedBuilds
	if err = json.NewDecoder(resp.Body).Decode(&all); err != nil {
		log.Printf("[%s][ERROR] Couldn't decode response : %v\n", s.Name, err)
		return
	}
	pall := make(Builds, len(all))
	for i, b := range all {
		pall[i] = b.Parse()
	}
	s.LastBuilds = pall
}

// FetchCommits fetches the last commits associated to the repository
func (s *Service) FetchCommits() {
	clog := logrus.WithFields(logrus.Fields{"action": "commits", "service": s.Name})
	u := strings.Split(s.Repo.URL, "/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", u[len(u)-2], u[len(u)-1])
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clog.WithError(err).Error("Couldn't create request")
		return
	}
	if conf.C.GithubOAuthToken != "" {
		req.Header.Add("Authorization", "token "+conf.C.GithubOAuthToken)
	}
	res, err := client.Do(req)
	if err != nil {
		clog.WithError(err).Warn("Couldn't perform request")
		return
	}
	defer res.Body.Close() //nolint:errcheck
	if res.StatusCode != http.StatusOK {
		clog.WithField("code", res.StatusCode).Warn("Couldn't retrieve commits")
		return
	}
	var all Commits
	if err = json.NewDecoder(res.Body).Decode(&all); err != nil {
		clog.WithError(err).Error("Couldn't decode response")
		return
	}
	s.LastCommits = all
}

// FetchRepoInfos fetches the repository information
func (s *Service) FetchRepoInfos() {
	clog := logrus.WithFields(logrus.Fields{"action": "repo", "service": s.Name})
	u := strings.Split(s.Repo.URL, "/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", u[len(u)-2], u[len(u)-1])
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clog.WithError(err).Error("Couldn't create request")
		return
	}
	if conf.C.GithubOAuthToken != "" {
		req.Header.Add("Authorization", "token "+conf.C.GithubOAuthToken)
	}
	res, err := client.Do(req)
	if err != nil {
		clog.WithError(err).Error("Couldn't perform request")
		return
	}
	defer res.Body.Close() //nolint:errcheck
	var repo GHRepo
	if err = json.NewDecoder(res.Body).Decode(&repo); err != nil {
		clog.WithError(err).Error("Couldn't decode response")
		return
	}
	s.Repo.Stars = repo.StargazersCount
	s.Repo.Forks = repo.ForksCount
	s.Repo.Watchers = repo.SubscribersCount
	s.Repo.Description = repo.Description
}

// Services represents a list of services
type Services []*Service

// Monitor allows to monitor Services every interval delay
func (ss Services) Monitor() {
	for _, s := range ss {
		if s.URL != "" {
			s.FetchStatus()
		}
	}
	for _, s := range ss {
		if s.CI != nil {
			go s.FetchBuilds()
		}
		if s.Repo != nil {
			go s.FetchCommits()
			go s.FetchRepoInfos()
		}
	}

	rtc := time.NewTicker(conf.C.RepoInterval)
	stc := time.NewTicker(conf.C.ServiceInterval)
	for {
		select {
		case <-rtc.C:
			logrus.WithField("type", "repo").Debug("Started background routine")
			for _, s := range ss {
				if s.CI != nil {
					go s.FetchBuilds()
				}
				if s.Repo != nil {
					go s.FetchCommits()
					go s.FetchRepoInfos()
				}
			}
		case <-stc.C:
			logrus.WithField("type", "status").Debug("Started background routine")
			for _, s := range ss {
				if s.URL != "" {
					s.FetchStatus()
				}
			}
		}
	}
}
