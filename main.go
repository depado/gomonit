package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/conf"
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
	return r
}

func main() {
	var err error

	if err = conf.Load("conf.yml"); err != nil {
		logrus.WithError(err).Fatal("Couldn't load configuration")
	}

	if err = models.InitializeServices(); err != nil {
		logrus.WithError(err).Fatal("Couldn't initialize services")
	}

	// Starting monitoring of services
	go models.All.Monitor()

	// Gin initialization
	if !conf.C.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set router
	r := SetupRouter()
	logrus.WithFields(logrus.Fields{
		"port": conf.C.Server.Port,
		"host": conf.C.Server.Host,
	}).Info("Starting server")
	if err = r.Run(fmt.Sprintf("%s:%d", conf.C.Server.Host, conf.C.Server.Port)); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
