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
	Branch string `json:"branch"`
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

		circleCIRequest := circleCIRequestType{
			"release",
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(circleCIRequest)

		headers := map[string][]string{
			"Content-Type":           []string{"application/json"},
			"Accept":                 []string{"application/json"},
			"x-attribution-login":    []string{"string"},
			"x-attribution-actor-id": []string{"string"},
		}

		encodedProjectSlug := url.QueryEscape("github/rivernews/iriversland2-kubernetes")
		var urlBuilder strings.Builder
		urlBuilder.WriteString("https://circleci.com/api/v2/project/")
		urlBuilder.WriteString(encodedProjectSlug)
		urlBuilder.WriteString("/pipeline")
		log.Printf("requesting circle ci at %s", urlBuilder.String())

		circleCiRequestURL, _ := url.Parse(urlBuilder.String())
		circleCiRequestURL.RawQuery = params.Encode()

		req, err := http.NewRequest("POST", circleCiRequestURL.String(), buf)
		req.Header = headers
		client := &http.Client{}
		res, err := client.Do(req)

		var slackMessage strings.Builder
		slackMessage.WriteString("K8s header service triggered circle ci job, response:\n```\n")
		bytesContent, _ := ioutil.ReadAll(res.Body)
		slackMessage.WriteString(string(bytesContent))
		slackMessage.WriteString("\n```\nAny error:\n```\n")
		if err != nil {
			slackMessage.WriteString("🔴 ")
			slackMessage.WriteString(err.Error())
		} else {
			slackMessage.WriteString("🟢 No error")
		}
		slackMessage.WriteString("\n```\n")
		c.JSON(http.StatusOK, gin.H{
			"text": slackMessage.String(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}
