package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type slackRequestType struct {
	Token string `json:"token" binding:"required"`
}

var requestFromSlackToken, requestFromSlackTokenExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")

// in order to export this function you need to capitalize it
// https://tour.golang.org/basics/3
func SlackController(c *gin.Context) {
	log.Println("in slack controller")

	if !requestFromSlackTokenExists {
		log.Panic(errors.New("slack token not configured"))
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "slack auth token not set",
		})
		return
	}

	slackRequest := slackRequestType{}
	if err := c.ShouldBindBodyWith(&slackRequest, binding.JSON); err != nil {
		log.Printf("Cannot parse slack request, ignored: %s", err)
		requestToken, exists := c.GetPostForm("token")
		log.Printf("GetPostForm(): %s, %s", requestToken, exists)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf("slack token received!")
	if requestFromSlackToken == slackRequest.Token {
		c.JSON(http.StatusOK, gin.H{
			"text": "K8s header service response:\n```received!```\n",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}
