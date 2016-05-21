package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/admin"
	"github.com/Depado/gomonit/auth"
	"github.com/Depado/gomonit/models"
	"github.com/Depado/gomonit/views"
)

// APIVersion is the current API version.
const APIVersion = "1"

// SetupRouter sets up the router and its routes as well as the templates and
// static files routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")

	// Main view
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
	return r
}

func main() {
	var err error

	// Configuration parsing and services initialization
	if err = models.ParseConfiguration("conf.yml"); err != nil {
		log.Fatal(err)
	}
	// Starting monitoring of services
	go models.All.Monitor()

	// Gin initialization
	if !models.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := SetupRouter()
	router.Run(fmt.Sprintf("%s:%s", models.C.Host, models.C.Port))
}
