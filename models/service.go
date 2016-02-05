package models

import (
	"log"
	"net/http"
	"time"
)

// Service is a single service
type Service struct {
	Name     string
	URL      string
	ShortURL string
	Host     string
	Last     string
	RespTime time.Duration
	Status   int
	Up       bool
	Icon     string
}

// Check updates the status of the Host
func (s *Service) Check(client *http.Client) {
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

// Services represents a list of services
type Services []*Service

// ServiceForm is the struct representing a Service (to add, or modify)
type ServiceForm struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	ShortURL string `form:"shorturl" binding:"required"`
}
