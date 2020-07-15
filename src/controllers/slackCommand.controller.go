package controllers

import (
	"log"
	"net/http"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/queue"
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
	parsedSlackRequest := types.SlackRequestType{}
	if err := c.ShouldBind(&parsedSlackRequest); err != nil {
		log.Printf("Cannot parse slack request, ignored: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if utilities.RequestFromSlackTokenCredential == parsedSlackRequest.Token {
		if parsedSlackRequest.TriggerWord == "slk" {
			utilities.TravisCITriggerSLKHelper(parsedSlackRequest)
		} else if parsedSlackRequest.TriggerWord == "kkk" || parsedSlackRequest.TriggerWord == "ddd" {
			utilities.CircleCITriggerK8sClusterHelper(parsedSlackRequest)
		} else if parsedSlackRequest.TriggerWord == "guide" {
			queue.HandleJobQueueRequest()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": utilities.BuildString(
					"no such command: ",
					parsedSlackRequest.TriggerWord,
				),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": "OK",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "auth failed",
	})
}
