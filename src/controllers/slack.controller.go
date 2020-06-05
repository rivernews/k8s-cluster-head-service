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
	if !requestFromSlackTokenExists {
		log.Panic(errors.New("slack token not configured"))
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "slack auth token not set",
		})
	}

	slackRequest := slackRequestType{}
	if err := c.ShouldBindBodyWith(&slackRequest, binding.JSON); err == nil {
		log.Printf("slack token: %s", slackRequest.Token)

		if requestFromSlackToken == slackRequest.Token {
			c.JSON(http.StatusOK, gin.H{
				"text": "K8s header service response:\n```received!```\n",
			})
		}
	}

	log.Println("cannot parse slack request, ignored")
}
