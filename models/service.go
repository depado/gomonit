package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Service is a single service
type Service struct {
	Name            string
	URL             string
	ShortURL        string
	RepoURL         string
	Host            string
	BuildAPI        string
	BuildURL        string
	CurrentBuildURL string
	Last            string
	RespTime        time.Duration
	Status          int
	Up              bool
	Icon            string
	LastBuild       string
	LastBuildTime   time.Duration
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
		s.Up = false
		log.Printf("[%s][ERROR] While building request : %v\n", s.Name, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		s.Up = false
		log.Printf("[%s][ERROR] While requesting : %v\n", s.Name, err)
		return
	}
	defer resp.Body.Close()
	s.Status = resp.StatusCode
	s.Up = s.Status == 200
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
	var all Builds
	if err = json.NewDecoder(resp.Body).Decode(&all); err != nil {
		log.Printf("[%s][ERROR] Couldn't decode response : %v\n", s.Name, err)
		return
	}
	if len(all) > 0 {
		s.LastBuild = all[0].Status
		s.CurrentBuildURL = fmt.Sprintf("%s%v", s.BuildURL, all[0].Number)
	}
}

// Check updates the status of the Host
func (s *Service) Check(client *http.Client) {
	go s.CheckStatus(client)
	if s.BuildAPI != "" {
		go s.CheckBuild(client)
	}
}

// Services represents a list of services
type Services []*Service

// ServiceForm is the struct representing a Service (to add, or modify)
type ServiceForm struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	ShortURL string `form:"shorturl" binding:"required"`
}
