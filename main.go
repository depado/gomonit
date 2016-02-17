package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/admin"
	"github.com/Depado/gomonit/auth"
	"github.com/Depado/gomonit/configuration"
	"github.com/Depado/gomonit/models"
)

var all models.Services

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"all": all,
	})
}

func main() {
	var err error
	if err = configuration.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}
	cnf := configuration.C
	if all, err = cnf.Parse(); err != nil {
		log.Fatal(err)
	}
	go all.Monitor(cnf.UpdateInterval)
	if !cnf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
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
	r.Run(cnf.Listen)
}
