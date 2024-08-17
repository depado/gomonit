package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/depado/gomonit/models"
)

// Index is the main route
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"all": models.All,
	})
}
