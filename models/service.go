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

// Service is a single service
type Service struct {
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	ShortURL        string        `json:"short_url"`
	RepoURL         string        `json:"repo_url"`
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

// CheckStatus checks if the service is running
func (s *Service) CheckStatus(client *http.Client) {
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

// CheckBuild checks the last build
func (s *Service) CheckBuild(client *http.Client) {
	req, err := http.NewRequest("GET", s.BuildAPI, nil)
	if err != nil {
		log.Printf("[%s][ERROR] While building request for build : %v\n", s.Name, err)
		return
	}
	resp, err := client.Do(req)
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

// Check updates the status of the Host
func (s *Service) Check() {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	if s.URL != "" {
		go s.CheckStatus(client)
	}
	if s.BuildAPI != "" {
		go s.CheckBuild(client)
	}
	if s.RepoURL != "" {
		go s.FetchCommits(client)
	}
}

// FetchCommits fetches the last commits associated to the repository
func (s *Service) FetchCommits(client *http.Client) {
	u := strings.Split(s.RepoURL, "/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", u[len(u)-2], u[len(u)-1])
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[%s][ERROR][COMMITS] While building request : %v\n", s.Name, err)
		return
	}
	resp, err := client.Do(req)
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
