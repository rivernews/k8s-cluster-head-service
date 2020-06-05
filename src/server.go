package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Default page
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// Slack routes
	slackRoutes := router.Group("/slack")
	slackRoutes.POST("provision", controllers.SlackController)

	// Protected routes

	accounts := gin.Accounts{}
	credential, exists := os.LookupEnv("CIRCLECI_TOKEN")
	if exists && credential != "" {
		accounts = gin.Accounts{
			"admin": credential,
		}
	}
	authorizedRoutes := router.Group("/protected", gin.BasicAuth(accounts))
	// endpoint /protected/
	authorizedRoutes.GET("/", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	})

	var address strings.Builder
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3010"
	}
	address.WriteString(":")
	address.WriteString(port)

	// Listen and serve on 0.0.0.0:8080
	router.Run(address.String())
}
