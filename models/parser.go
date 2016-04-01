package models

import (
	"fmt"
	"strings"

	"github.com/Depado/gomonit/conf"
)

// ParseConf parses the configuration and returns the appropriate Services
func ParseConf() error {
	All = make(Services, len(conf.C.Services))
	conf.C.CIURL = strings.TrimSuffix(conf.C.CIURL, "/")
	for i, s := range conf.C.Services {
		if s.CIType != "" && s.CIType != "drone" {
			return fmt.Errorf("Unable to use %s as CI, currently only drone is supported.", s.CIType)
		}
		if s.RepoType != "" && s.RepoType != "github" {
			return fmt.Errorf("Unable to use %s as repository, currently only github is supported.", s.RepoType)
		}
		repoURL := ""
		if s.Repo != "" {
			repoURL = fmt.Sprintf("https://github.com/%s", s.Repo)
		}
		buildAPI := ""
		buildURL := ""
		if s.CIType != "" {
			buildAPI = fmt.Sprintf("%s/api/repos/%s/builds", conf.C.CIURL, s.Repo)
			buildURL = fmt.Sprintf("%s/%s", conf.C.CIURL, s.Repo)
		}
		short := strings.TrimPrefix(s.URL, "http://")
		short = strings.TrimPrefix(short, "https://")
		All[i] = &Service{
			Name:     s.Name,
			URL:      s.URL,
			ShortURL: short,
			Host:     s.Host,
			BuildAPI: buildAPI,
			BuildURL: buildURL,
			RepoURL:  repoURL,
			Icon:     "/static/custom/" + s.Icon,
			Own:      s.Own,
		}
	}
	return nil
}
