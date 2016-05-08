// main.go (web-server)

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var auth = authService{Base: "http://localhost:8001"}

func main() {
	gin.SetMode(gin.ReleaseMode)
	s := gin.Default()

	s.POST("/login", login)
	s.POST("/logout", logout)

	s.Run(":8000")
}

// Handler for the login request
func login(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	if response := auth.Login(username, password); response.Token != "" {
		// If authentication succeeds respond with an HTTP success
		// status and include the token in the response
		c.JSON(http.StatusOK, response)
	} else {
		// Respond with an HTTP error if authentication fails
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// Handler for the logout request
func logout(c *gin.Context) {
	// Obtain the POSTed username and token values
	username := c.PostForm("username")
	token := c.PostForm("token")

	if auth.Logout(username, token) {
		// Respond with an HTTP success status
		c.JSON(http.StatusOK, nil)
	} else {
		// Respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
