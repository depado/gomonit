package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/admin"
	"github.com/Depado/gomonit/auth"
	"github.com/Depado/gomonit/configuration"
	"github.com/Depado/gomonit/models"
)

var all models.Services

func periodicHostUpdate() {
	tc := time.NewTicker(30 * time.Minute)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	for {
		for _, host := range all {
			go host.Check(client)
		}
		<-tc.C
	}
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"all": all,
	})
}

func main() {
	if err := configuration.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}
	cnf := configuration.C
	all = make(models.Services, len(cnf.Services))
	for i, s := range cnf.Services {
		all[i] = &models.Service{
			Name:     s.Name,
			URL:      s.URL,
			ShortURL: s.ShortURL,
			Host:     s.Host,
			Up:       false,
			Icon:     "/static/custom/" + s.Icon,
		}
	}
	go periodicHostUpdate()
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")

	r.GET("/", index)

	r.GET("/login", auth.Login)
	r.POST("/login", auth.PostLogin)

	ar := r.Group("/admin")
	{
		ar.GET("/", admin.Root)
		ar.GET("/hosts", admin.Hosts)
		ar.GET("/hosts/new", admin.NewHost)
		ar.POST("/hosts/new", admin.PostNewHost)
	}
	r.Run(fmt.Sprintf("127.0.0.1:%d", cnf.Port))
}
