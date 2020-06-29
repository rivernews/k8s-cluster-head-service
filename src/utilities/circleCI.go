package utilities

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"
)

var circleCIAPIBaseURL = "https://circleci.com/api/v2"

var circleCIHeaders = map[string][]string{
	"Content-Type":           {"application/json"},
	"Accept":                 {"application/json"},
	"x-attribution-login":    {"string"},
	"x-attribution-actor-id": {"string"},
}

func CircleCITriggerK8sClusterHelper(parsedSlackRequest types.SlackRequestType) {
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

	// prepare url static path parameter
	encodedProjectSlug := url.QueryEscape("github/rivernews/iriversland2-kubernetes")

	// generate api call url and assign static path parameter
	// var urlBuilder strings.Builder
	urlBuilder := StringBuilder([]string{
		circleCIAPIBaseURL,
		"/project/",
		encodedProjectSlug,
		"/pipeline",
	})
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
	_, fetchResultMessage, _ := Fetch(FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: circleCIHeaders,
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

// FetchCircleCIBuildStatus - returns the status string of the build given a pipeline uuid.
func FetchCircleCIBuildStatus(pipelineID string) (string, error) {
	fetchURL := BuildString([]string{
		circleCIAPIBaseURL, "/pipeline/", pipelineID, "/workflow",
	})

	responseBytes, _, fetchErr := Fetch(FetchOption{
		Method:  "GET",
		URL:     fetchURL,
		Headers: circleCIHeaders,
		QueryParams: map[string]string{
			"circle-token": CircleCiToken,
		},
		DisableHumanMessage: true,
	})
	if fetchErr != nil {
		return "", fetchErr
	}

	var pipelineWorkflows types.CircleCIPipelineWorkflowListResponseType
	unmarshalJSONErr := json.Unmarshal(responseBytes, &pipelineWorkflows)
	if unmarshalJSONErr != nil {
		return "", unmarshalJSONErr
	}

	if len(pipelineWorkflows.Workflows) < 1 {
		return "", errors.New("No workflow exists for pipeline " + pipelineID)
	}

	return pipelineWorkflows.Workflows[0].Status, nil
}

// CircleCIWaitTillWorkflowFinish - wait till the build finalizes,
// when finalized, return the final status string
func CircleCIWaitTillWorkflowFinish(pipelineID string) (string, error) {
	MaxPollingCount := 12 * 20
	pollingCount := 0

	// while look polling
	status := "on_hold"
	var fetchErr error = nil
	for (status == "running" || status == "on_hold") && pollingCount <= MaxPollingCount {
		status, fetchErr = FetchCircleCIBuildStatus(pipelineID)

		time.Sleep(5 * time.Second)
		pollingCount++
	}

	if status != "running" && status != "on_hold" && fetchErr == nil {
		return status, fetchErr
	}

	return "", errors.New(
		BuildString([]string{
			"Timed out while waiting for CircleCI workflow to finish for pipeline ", pipelineID, "\n",
			"Any fetch error? ", fetchErr.Error(),
		}),
	)
}
