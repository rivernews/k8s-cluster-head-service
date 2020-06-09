package controllers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"
	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"

	"github.com/gin-gonic/gin"
)

type slackRequestType struct {
	Token       string `form:"token"`
	TriggerWord string `form:"trigger_word"`
	Text        string `form:"text"`
}

// SlackController port slack command to circle CI API
//
// Projec status
// https://app.circleci.com/pipelines/github/rivernews/iriversland2-kubernetes
//
// API doc
// https://circleci.com/docs/api/v2/?shell#trigger-a-new-pipeline
func SlackController(c *gin.Context) {
	log.Println("in slack controller")

	if !utilities.RequestFromSlackTokenCredentialExists {
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

	if utilities.RequestFromSlackTokenCredential == parsedSlackRequest.Token {
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
	// parse slack command
	fullCommand := strings.TrimSpace(parsedSlackRequest.Text)
	fullCommand = strings.ToLower(fullCommand)
	commandTokens := strings.Split(fullCommand, ":")
	// parse dropletSize
	dropletSize := utilities.KubernetesClusterDefaultDropletSize
	if len(commandTokens) > 1 {
		if commandTokens[1] == "m" {
			dropletSize = utilities.MediumDroplet
		} else if commandTokens[1] == "s" {
			dropletSize = utilities.SmallDroplet
		}
	}

	// prepare post data
	branch := "master"
	if parsedSlackRequest.TriggerWord == "kkk" {
		branch = "release"
	} else if parsedSlackRequest.TriggerWord == "ddd" {
		branch = "destroy-release"
	}
	envVars := map[string]string{}
	envVars["TF_VAR_droplet_size"] = dropletSize

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

	// make request
	var responseMessage strings.Builder
	if branch == "release" {
		responseMessage.WriteString("Provisioning kubernetes requested.\n")
	} else if branch == "destroy-release" {
		responseMessage.WriteString("Destroying kubernetes requested.\n")
	} else {
		responseMessage.WriteString("Verify kubernetes requested.\n")
	}
	fetchResultMessage := utilities.Fetch(utilities.FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: headers,
		QueryParams: map[string]string{
			"circle-token": utilities.CircleCiToken,
		},
		PostData: types.CircleCIRequestType{
			Branch:  branch,
			EnvVars: envVars,
		},
	})
	responseMessage.WriteString(fetchResultMessage)

	c.JSON(http.StatusOK, gin.H{
		"text": responseMessage,
	})

	return
}

var travisAPIBaseURL = "https://api.travis-ci.com"

func travisCITriggerSLKHelper(c *gin.Context, parsedSlackRequest slackRequestType) {
	encodedProjectSlug := url.QueryEscape("rivernews/slack-middleware-server")

	// build url
	var urlBuilder strings.Builder
	urlBuilder.WriteString(travisAPIBaseURL)
	// endpoint
	urlBuilder.WriteString("/repo/")
	urlBuilder.WriteString(encodedProjectSlug)
	urlBuilder.WriteString("/requests")

	fetchedMessage := utilities.Fetch(utilities.FetchOption{
		Method: "POST",
		URL:    urlBuilder.String(),
		Headers: map[string][]string{
			"Content-Type":       {"application/json"},
			"Accept":             {"application/json"},
			"Travis-API-Version": {"3"},
			"Authorization":      {"token " + utilities.TravisCIToken},
		},
		PostData: map[string]string{
			"branch": "release",
		},
	})

	var respondSlackMessage strings.Builder
	respondSlackMessage.WriteString("Provision SLK requested.\n")
	respondSlackMessage.WriteString(fetchedMessage)

	c.JSON(http.StatusOK, gin.H{
		"text": respondSlackMessage,
	})
}
