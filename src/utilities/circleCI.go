package utilities

import (
	"log"
	"net/url"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"

	"github.com/gin-gonic/gin"
)

func CircleCITriggerK8sClusterHelper(c *gin.Context, parsedSlackRequest types.SlackRequestType) {
	// parse slack command
	fullCommand := strings.TrimSpace(parsedSlackRequest.Text)
	fullCommand = strings.ToLower(fullCommand)
	commandTokens := strings.Split(fullCommand, ":")
	// parse dropletSize - default is large size
	dropletSize := LargeDroplet
	if len(commandTokens) > 1 {
		if commandTokens[1] == "m" {
			dropletSize = MediumDroplet
		} else if commandTokens[1] == "s" {
			dropletSize = SmallDroplet
		} else if commandTokens[1] == "l" {
			dropletSize = LargeDroplet
		}
	}

	// prepare post data
	branch := "master"
	if parsedSlackRequest.TriggerWord == "kkk" {
		branch = "release"
	} else if parsedSlackRequest.TriggerWord == "ddd" {
		branch = "destroy-release"
	}

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
		responseMessage.WriteString("Droplet size: `")
		responseMessage.WriteString(dropletSize)
		responseMessage.WriteString("`\n")
	} else if branch == "destroy-release" {
		responseMessage.WriteString("Destroying kubernetes requested.\n")
	} else {
		responseMessage.WriteString("Verify kubernetes requested.\n")
	}
	fetchResultMessage := Fetch(FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: headers,
		QueryParams: map[string]string{
			"circle-token": CircleCiToken,
		},
		PostData: types.CircleCIRequestType{
			Branch: branch,
			Parameters: types.CircleCIKubernetesClusterProjectPipelineParameters{
				DropletSize: dropletSize,
			},
		},
	})
	responseMessage.WriteString(fetchResultMessage)
	responseMessage.WriteString("<https://app.circleci.com/pipelines/github/rivernews/iriversland2-kubernetes|Check out the pipeline> in CircleCI dashboard.\n")
	
	SendSlackMessage(responseMessage.String())

	return
}
