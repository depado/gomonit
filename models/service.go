package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// All represents all the services
var All Services

// Service is a single service
type Service struct {
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	ShortURL        string        `json:"short_url"`
	Description     string        `json:"description"`
	RepoURL         string        `json:"repo_url"`
	RepoStars       int           `json:"repo_stars"`
	RepoForks       int           `json:"repo_forks"`
	RepoWatchers    int           `json:"repo_watchers"`
	Host            string        `json:"host"`
	BuildAPI        string        `json:"build_api"`
	BuildURL        string        `json:"build_url"`
	CurrentBuildURL string        `json:"current_build_url"`
	Last            string        `json:"last"`
	RespTime        time.Duration `json:"resp_time"`
	Status          int           `json:"status"`
	Icon            string        `json:"icon"`
	LastBuilds      Builds        `json:"last_builds"`
	LastCommits     Commits       `json:"last_commits"`
	Own             bool          `json:"own"`
}

// FetchStatus checks if the service is running
func (s *Service) FetchStatus(client *http.Client) {
	start := time.Now()
	s.Last = start.Format("2006/01/02 15:04:05")
	defer func() {
		d := time.Since(start)
		s.RespTime = d - (d % time.Millisecond)
	}()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		log.Printf("[%s][ERROR] While building request : %v\n", s.Name, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%s][ERROR] While requesting : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close()
	s.Status = resp.StatusCode
}

// FetchBuilds checks the last build
func (s *Service) FetchBuilds(client *http.Client) {
	resp, err := http.Get(s.BuildAPI)
	if err != nil {
		log.Printf("[%s][ERROR] While requesting build status : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close()
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
func (s *Service) FetchCommits(client *http.Client) {
	u := strings.Split(s.RepoURL, "/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", u[len(u)-2], u[len(u)-1])
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[%s][ERROR][COMMITS] While requesting : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close()
	var all Commits
	if err = json.NewDecoder(resp.Body).Decode(&all); err != nil {
		log.Printf("[%s][ERROR][COMMITS] Couldn't decode response : %v\n", s.Name, err)
		return
	}
	s.LastCommits = all
}

// FetchRepoInfos fetches the repository information
func (s *Service) FetchRepoInfos(client *http.Client) {
	u := strings.Split(s.RepoURL, "/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", u[len(u)-2], u[len(u)-1])
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[%s][ERROR][REPO] While requesting : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close()
	var repo GHRepo
	if err = json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		log.Printf("[%s][ERROR][REPO] Couldn't decode response : %v\n", s.Name, err)
		return
	}
	s.RepoStars = repo.StargazersCount
	s.RepoForks = repo.ForksCount
	s.RepoWatchers = repo.SubscribersCount
	s.Description = repo.Description
}

// Check updates the status of the Host
func (s *Service) Check() {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	if s.URL != "" {
		s.FetchStatus(client)
	}
	if s.BuildAPI != "" {
		go s.FetchBuilds(client)
	}
	if s.RepoURL != "" {
		go s.FetchCommits(client)
		go s.FetchRepoInfos(client)
	}
}

// Services represents a list of services
type Services []*Service

// Monitor allows to monitor Services every interval delay
func (ss Services) Monitor(interval time.Duration) {
	tc := time.NewTicker(interval)
	for {
		for _, s := range ss {
			go s.Check()
		}
		<-tc.C
	}
}

// ServiceForm is the struct representing a Service (to add, or modify)
type ServiceForm struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	ShortURL string `form:"shorturl" binding:"required"`
}
