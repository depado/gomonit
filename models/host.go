package models

import (
	"log"
	"net/http"
	"time"
)

// Host is a single host
type Host struct {
	Name     string
	URL      string
	ShortURL string
	Last     string
	RespTime time.Duration
	Status   int
	Up       bool
	Icon     string
}

// Check updates the status of the Host
func (h *Host) Check(client *http.Client) {
	start := time.Now()
	h.Last = start.Format("2006/01/02 15:04:05")
	defer func() {
		d := time.Since(start)
		h.RespTime = d - (d % time.Millisecond)
	}()
	req, err := http.NewRequest("GET", h.URL, nil)
	if err != nil {
		h.Up = false
		log.Printf("[%s][ERROR] While building request : %v\n", h.Name, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		h.Up = false
		log.Printf("[%s][ERROR] While requesting : %v\n", h.Name, err)
		return
	}
	h.Status = resp.StatusCode
	h.Up = h.Status == 200
}

// Hosts are the hosts
type Hosts []*Host

// HostForm is the struct representing a host (to add, or modify)
type HostForm struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	ShortURL string `form:"shorturl" binding:"required"`
}
