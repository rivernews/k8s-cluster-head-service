package controllers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"

	"github.com/gin-gonic/gin"
)

type slackRequestType struct {
	Token       string `form:"token"`
	TriggerWord string `form:"trigger_word"`
}

var requestFromSlackTokenCredential, requestFromSlackTokenCredentialExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")
var circleCiToken, _ = os.LookupEnv("CIRCLECI_TOKEN")

// SlackController port slack command to circle CI API
//
// Projec status
// https://app.circleci.com/pipelines/github/rivernews/iriversland2-kubernetes
//
// API doc
// https://circleci.com/docs/api/v2/?shell#trigger-a-new-pipeline
func SlackController(c *gin.Context) {
	log.Println("in slack controller")

	if !requestFromSlackTokenCredentialExists {
		log.Panic(errors.New("slack token not configured"))
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "slack auth token not set",
		})
		return
	}

	parsedSlackRequest := slackRequestType{}
	if err := c.ShouldBind(&parsedSlackRequest); err != nil {
		log.Printf("Cannot parse slack request, ignored: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if requestFromSlackTokenCredential == parsedSlackRequest.Token {
		// TODO: cancel job
		if parsedSlackRequest.TriggerWord == "ppp" {
			// TODO: poll job status
		} else if parsedSlackRequest.TriggerWord == "slk" {
			travisCITriggerSLKHelper(c, parsedSlackRequest)
		} else {
			circleCITriggerK8sClusterHelper(c, parsedSlackRequest)
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}

func circleCITriggerK8sClusterHelper(c *gin.Context, parsedSlackRequest slackRequestType) {
	// prepare credentials via querystring
	// params := url.Values{}
	// params.Add("circle-token", circleCiToken)

	// prepare post data
	branch := "master"
	if parsedSlackRequest.TriggerWord == "kkk" {
		branch = "release"
	} else if parsedSlackRequest.TriggerWord == "ddd" {
		branch = "destroy-release"
	}
	// circleCIRequest := types.CircleCIRequestType{Branch: branch}
	// circleCIPOSTRequestFormBuf := new(bytes.Buffer)
	// json.NewEncoder(circleCIPOSTRequestFormBuf).Encode(circleCIRequest)

	// prepare headers
	headers := map[string][]string{
		"Content-Type":           {"application/json"},
		"Accept":                 {"application/json"},
		"x-attribution-login":    {"string"},
		"x-attribution-actor-id": {"string"},
	}

	// prepare url static path parameter
	encodedProjectSlug := url.QueryEscape("github/rivernews/iriversland2-kubernetes")

	// generate api call url and assign static path parameter
	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://circleci.com/api/v2/project/")
	urlBuilder.WriteString(encodedProjectSlug)
	urlBuilder.WriteString("/pipeline")
	log.Printf("requesting circle ci at %s", urlBuilder.String())

	// add credentials by querystring
	// circleCiRequestURL, _ := url.Parse(urlBuilder.String())
	// circleCiRequestURL.RawQuery = params.Encode()

	// append request config and make request
	// req, err := http.NewRequest("POST", circleCiRequestURL.String(), circleCIPOSTRequestFormBuf)
	// req.Header = headers
	// client := &http.Client{}
	// res, err := client.Do(req)

	// log response
	// var slackMessage strings.Builder
	// slackMessage.WriteString("K8s header service triggered circle ci job, response:\n```\n")
	// bytesContent, _ := ioutil.ReadAll(res.Body)
	// slackMessage.WriteString(string(bytesContent))
	// slackMessage.WriteString("\n```\nAny error:\n```\n")
	// if err != nil {
	// 	slackMessage.WriteString("🔴 ")
	// 	slackMessage.WriteString(err.Error())
	// } else {
	// 	slackMessage.WriteString("🟢 No error")
	// }
	// slackMessage.WriteString("\n```\n<")
	// projectDashboardURL := "https://app.circleci.com/pipelines/github/rivernews/iriversland2-kubernetes"
	// slackMessage.WriteString(projectDashboardURL)
	// slackMessage.WriteString("|Project dashboard>.")

	slackMessage := ""
	utilities.Fetch(utilities.FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: headers,
		QueryParams: map[string]string{
			"circle-token": circleCiToken,
		},
		PostData: map[string]string{
			"branch": branch,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"text": slackMessage,
	})

	return
}

func travisCITriggerSLKHelper(c *gin.Context, parsedSlackRequest slackRequestType) {

}
