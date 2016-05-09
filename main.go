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
	s.GET("/logout", logout)
	s.GET("/protected-content", serveProtectedContent)

	s.Run(":8000")
}

// Handler for the login request
func login(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	if response := auth.Login(username, password); response.Token != "" {
		// If authentication succeeds set the cookies and
		// respond with an HTTP success
		// status and include the token in the response
		c.SetCookie("username", username, 3600, "", "", false, true)
		c.SetCookie("token", response.Token, 3600, "", "", false, true)

		c.JSON(http.StatusOK, response)
	} else {
		// Respond with an HTTP error if authentication fails
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// Handler for the logout request
func logout(c *gin.Context) {
	// Obtain the username and token from the cookies
	username, err1 := c.Cookie("username")
	token, err2 := c.Cookie("token")

	if err1 == nil && err2 == nil && auth.Logout(username, token) {
		// Clear the cookies and
		// respond with an HTTP success status
		c.SetCookie("username", "", -1, "", "", false, true)
		c.SetCookie("token", "", -1, "", "", false, true)

		c.JSON(http.StatusOK, nil)
	} else {
		// Respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

// Handler to serve the protected content
func serveProtectedContent(c *gin.Context) {
	// Obtain the username and token from the cookies
	username, err1 := c.Cookie("username")
	token, err2 := c.Cookie("token")

	if err1 == nil && err2 == nil && auth.Authenticate(username, token) {
		// Respond with an HTTP success status & include the
		// content in the response

		c.JSON(http.StatusOK, gin.H{"content": "This should be visible to authenticated users only."})
	} else {
		// Respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
