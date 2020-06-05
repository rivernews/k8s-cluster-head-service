package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type slackRequestType struct {
	Token string `form:"token" json:"token"`
}

type circleCIRequestType struct {
	ProjectSlug string `json:"project-slug"`
	Branch      string `json:"branch"`
}

var requestFromSlackToken, requestFromSlackTokenExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")
var circleCiToken, _ = os.LookupEnv("CIRCLECI_TOKEN")

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
	if err := c.ShouldBind(&slackRequest); err != nil {
		log.Printf("Cannot parse slack request, ignored: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if requestFromSlackToken == slackRequest.Token {
		params := url.Values{}
		params.Add("circle-token", circleCiToken)
		circleCiRequestURL, _ := url.Parse("https://circleci.com/api/v2/project")
		circleCiRequestURL.RawQuery = params.Encode()

		circleCIRequest := circleCIRequestType{
			"github/rivernews/iriversland2-kubernetes",
			"release",
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(circleCIRequest)
		res, err := http.Post(circleCiRequestURL.String(), "application/json", buf)

		var slackMessage strings.Builder
		slackMessage.WriteString("K8s header service triggered circle ci job, response:\n```")
		bytesContent, _ := ioutil.ReadAll(res.Body)
		slackMessage.WriteString(string(bytesContent))
		slackMessage.WriteString("```\nAny error:\n```")
		slackMessage.WriteString(err.Error())
		slackMessage.WriteString("```\n")
		c.JSON(http.StatusOK, gin.H{
			"text": slackMessage.String(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}
