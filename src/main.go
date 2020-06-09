package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/controllers"
	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"

	"github.com/gin-gonic/gin"
)

func main() {
	if !checkAppConfigurationOK() {
		return
	}

	if !utilities.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Default page
	router.GET("/", func(c *gin.Context) {
		utilities.SendSlackMessage("haha")
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// Slack routes
	slackRoutes := router.Group("/slack")
	slackRoutes.POST("provision", controllers.SlackCommandController)

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
	log.Println("Start listening on " + address.String())
	router.Run(address.String())
}

func checkAppConfigurationOK() bool {
	log.Println("Checking app configuration...")

	if !utilities.RequestFromSlackTokenCredentialExists {
		log.Fatalln("Outgoing slack webhook token is not configured")
		return false
	}

	if !utilities.CircleCiTokenExists {
		log.Fatalln("CircleCI token is not configured")
		return false
	}

	if !utilities.TravisCITokenExists {
		log.Fatalln("TravisCI token is not configured")
		return false
	}

	if !utilities.SendSlackURLExists {
		log.Fatalln("Send slack URL is not configured")
		return false
	}

	return true
}
