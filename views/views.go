package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/gomonit/models"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"all": models.All,
	})
}
