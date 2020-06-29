package utilities

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"

	"github.com/gin-gonic/gin"
)

var travisAPIBaseURL = "https://api.travis-ci.com"
var travisCISLKEncodedProjectSlug = url.QueryEscape("rivernews/slack-middleware-server")
var travisCIAPIHeaders = map[string][]string{
	"Content-Type":       {"application/json"},
	"Accept":             {"application/json"},
	"Travis-API-Version": {"3"},
	"Authorization":      {"token " + TravisCIToken},
}

/*
	TravisCITriggerSLKHelper - triggers provisioning SLK deployment
*/
func TravisCITriggerSLKHelper(c *gin.Context, parsedSlackRequest types.SlackRequestType) {

	// build url
	var urlBuilder strings.Builder
	urlBuilder.WriteString(travisAPIBaseURL)
	// endpoint
	urlBuilder.WriteString("/repo/")
	urlBuilder.WriteString(travisCISLKEncodedProjectSlug)
	urlBuilder.WriteString("/requests")

	_, fetchedMessage, _ := Fetch(FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: travisCIAPIHeaders,
		PostData: map[string]string{
			"branch": "release",
		},
	})

	var respondSlackMessage strings.Builder
	respondSlackMessage.WriteString("Provision SLK requested.\n")
	respondSlackMessage.WriteString(fetchedMessage)

	SendSlackMessage(respondSlackMessage.String())

	return
}

func TravisCIWaitUntilBuildProvisioned(requestID string) (types.TravisCIBuildRequestResponseType, error) {
	// wait up to 5 minutes
	MaxPollingCount := 12 * 5
	pollingCount := 0

	travisCIRequestObject := types.TravisCIBuildRequestResponseType{
		Builds: []types.TravisCIBuild{},
	}

	fetchURL := BuildString([]string{
		travisAPIBaseURL, "/repo/", travisCISLKEncodedProjectSlug, "/request/", requestID,
	})

	for pollingCount <= MaxPollingCount {
		time.Sleep(5 * time.Second)
		responseBytes, _, fetchErr := Fetch(FetchOption{
			Method:              "GET",
			URL:                 fetchURL,
			Headers:             travisCIAPIHeaders,
			DisableHumanMessage: true,
		})
		pollingCount++

		if fetchErr != nil {
			log.Printf("Fetch error while waiting for travis build provisioned at request ID %s: %s", requestID, fetchErr)
			continue
		}

		unmarshalJSONErr := json.Unmarshal(responseBytes, &travisCIRequestObject)

		if unmarshalJSONErr != nil {
			log.Panicf("Unmarshal failed while waiting for travisCI build provisioned: %s", unmarshalJSONErr)
			continue
		}

		if len(travisCIRequestObject.Builds) > 0 {
			return travisCIRequestObject, nil
		}
	}

	return travisCIRequestObject, errors.New("Time out while waiting for travis build be provisioned for request ID " + requestID)
}

/*
	TravisCICheckBuildStutus

	Travis API doc for build:
	https://developer.travis-ci.com/resource/build#Build
*/
func TravisCICheckBuildStutus(buildID string) (string, error) {
	fetchURL := BuildString([]string{
		travisAPIBaseURL, "/build/", buildID,
	})

	responseBytes, _, fetchErr := Fetch(FetchOption{
		Method:              "GET",
		URL:                 fetchURL,
		Headers:             travisCIAPIHeaders,
		DisableHumanMessage: true,
	})
	if fetchErr != nil {
		return "", fetchErr
	}

	var travisCIBuild types.TravisCIBuild
	unmarshalJSONErr := json.Unmarshal(responseBytes, &travisCIBuild)
	if unmarshalJSONErr != nil {
		return "", unmarshalJSONErr
	}

	return travisCIBuild.State, nil
}
