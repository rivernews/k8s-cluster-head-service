package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"
	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"

	"github.com/gin-gonic/gin"
)

// SlackCommandController port slack command to circle CI API
//
// Projec status
// https://app.circleci.com/pipelines/github/rivernews/iriversland2-kubernetes
//
// API doc
// https://circleci.com/docs/api/v2/?shell#trigger-a-new-pipeline
//
// Pipeline parameter doc
// https://github.com/CircleCI-Public/api-preview-docs/blob/master/docs/pipeline-parameters.md
func SlackCommandController(c *gin.Context) {
	log.Println("in slack controller")

	if !utilities.RequestFromSlackTokenCredentialExists {
		log.Panic(errors.New("slack token not configured"))
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "slack auth token not set",
		})
		return
	}

	parsedSlackRequest := types.SlackRequestType{}
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
			utilities.TravisCITriggerSLKHelper(c, parsedSlackRequest)
		} else {
			utilities.CircleCITriggerK8sClusterHelper(c, parsedSlackRequest)
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}
