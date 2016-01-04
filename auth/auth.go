package auth

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/Depado/gomonit/models"
	"github.com/gin-gonic/gin"
)

// Login returns the HTML form to login
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

// PostLogin handles the post of the form
func PostLogin(c *gin.Context) {
	var lf models.LoginForm
	if c.Bind(&lf) == nil {
		log.Println("Login :", lf.Login)
		log.Println("Password :", lf.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(lf.Password), bcrypt.DefaultCost)
		if err == nil {
			log.Println(string(hashedPassword))
		}
	}
}
