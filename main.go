package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/admin"
	"github.com/Depado/gomonit/auth"
	"github.com/Depado/gomonit/conf"
	"github.com/Depado/gomonit/models"
	"github.com/Depado/gomonit/views"
)

func main() {
	var err error

	// Configuration parsing and services initialization
	if err = conf.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}
	if err = conf.C.Parse(); err != nil {
		log.Fatal(err)
	}
	// Starting monitoring of services
	go models.All.Monitor(conf.C.UpdateInterval)

	// Gin initialization
	if !conf.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")

	r.GET("/", views.Index)

	// API routes declaration
	api := r.Group("/api")
	{
		api.GET("/status", views.Status)
		api.GET("/dump/all", views.DumpAll)
		api.GET("/dump/own", views.DumpOwn)
	}

	// Authentication routes declaration
	r.GET("/login", auth.Login)
	r.POST("/login", auth.PostLogin)

	// Admin routes declaration
	ar := r.Group("/admin")
	{
		ar.GET("/", admin.Root)
		ar.GET("/hosts", admin.Hosts)
		ar.GET("/hosts/new", admin.NewHost)
		ar.POST("/hosts/new", admin.PostNewHost)
	}

	// Running
	r.Run(conf.C.Listen)
}
