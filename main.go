package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/admin"
	"github.com/Depado/gomonit/auth"
	"github.com/Depado/gomonit/configuration"
	"github.com/Depado/gomonit/models"
)

var all models.Hosts

func periodicHostUpdate() {
	tc := time.NewTicker(30 * time.Minute)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	for {
		for _, host := range all {
			go func(h *models.Host) {
				log.Println("Checking host", h.Name)
				h.Last = time.Now().Format("2006/01/02 15:04:05")
				req, err := http.NewRequest("GET", h.URL, nil)
				if err != nil {
					h.Up = false
					return
				}
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode != 200 {
					h.Up = false
					return
				}
				log.Println("Host", h.Name, "is UP")
				h.Up = true
			}(host)
			log.Println(&host)
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
	all = make(models.Hosts, len(cnf.Hosts))
	for i, h := range cnf.Hosts {
		all[i] = &models.Host{h.Name, h.URL, h.ShortURL, time.Now().Format("2006/01/02 15:04:05"), false, "/static/custom/" + h.Icon}
	}
	go periodicHostUpdate()
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

	r.Run(":8080")
}
