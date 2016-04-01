package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/models"
)

// Status gets only the status of all the services (HTTP status code)
func Status(c *gin.Context) {
	resp := gin.H{}
	for _, s := range models.All {
		resp[s.Name] = s.Status
	}
	c.JSON(200, resp)
}

// DumpAll dumps all the data and returns them as JSON
func DumpAll(c *gin.Context) {
	c.JSON(200, models.All)
}

// DumpOwn is the same as DumpAll but only for services marked as "own" in the
// configuration
func DumpOwn(c *gin.Context) {
	resp := models.Services{}
	for _, s := range models.All {
		if s.Own {
			resp = append(resp, s)
		}
	}
	c.JSON(http.StatusOK, resp)
}
