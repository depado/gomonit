package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/models"
)

func Status(c *gin.Context) {
	resp := gin.H{}
	for _, s := range models.All {
		resp[s.Name] = s.Status
	}
	c.JSON(200, resp)
}

func DumpAll(c *gin.Context) {
	c.JSON(200, models.All)
}

func DumpOwn(c *gin.Context) {
	resp := models.Services{}
	for _, s := range models.All {
		if s.Own {
			resp = append(resp, s)
		}
	}
	c.JSON(http.StatusOK, resp)
}
