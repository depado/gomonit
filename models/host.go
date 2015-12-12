package models

// Host is a single host
type Host struct {
	Name     string
	URL      string
	ShortURL string
	Last     string
	Up       bool
	Icon     string
}

// Hosts are the hosts
type Hosts []*Host

// HostForm is the struct representing a host (to add, or modify)
type HostForm struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	ShortURL string `form:"shorturl" binding:"required"`
}
