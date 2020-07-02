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
func TravisCITriggerSLKHelper(parsedSlackRequest types.SlackRequestType) (types.TravisCIRequestProvisionType, error) {

	// build url
	var urlBuilder strings.Builder
	urlBuilder.WriteString(travisAPIBaseURL)
	// endpoint
	urlBuilder.WriteString("/repo/")
	urlBuilder.WriteString(travisCISLKEncodedProjectSlug)
	urlBuilder.WriteString("/requests")

	var travisCIRequestProvision types.TravisCIRequestProvisionType

	_, fetchedMessage, fetchErr := Fetch(FetchOption{
		Method:  "POST",
		URL:     urlBuilder.String(),
		Headers: travisCIAPIHeaders,
		PostData: map[string]string{
			"branch": "release",
		},
		responseStore: travisCIRequestProvision,
	})

	var respondSlackMessage strings.Builder
	respondSlackMessage.WriteString("Provision SLK requested.\n")
	respondSlackMessage.WriteString(fetchedMessage)

	SendSlackMessage(respondSlackMessage.String())

	return travisCIRequestProvision, fetchErr
}

/*
travisCIWaitUntilBuildProvisioned - polls a request status,
and return the first (latest) build's state as soon as a build is provisioned.
*/
func travisCIWaitUntilBuildProvisioned(requestID string) (string, error) {
	// wait up to 5 minutes
	MaxPollingCount := 12 * 5
	pollingCount := 0

	travisCIRequestObject := types.TravisCIBuildRequestResponseType{
		Builds: []types.TravisCIBuild{},
	}

	fetchURL := BuildString(
		travisAPIBaseURL, "/repo/", travisCISLKEncodedProjectSlug, "/request/", requestID,
	)

	for pollingCount <= MaxPollingCount {
		if GetLogLevelValue() >= LogLevelTypes["INFO"] {
			log.Print(BuildString(
				"Polling ", requestID, " ...",
			))
		}

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
			return travisCIRequestObject.Builds[0].State, nil
		}
	}

	return "", errors.New("Time out while waiting for travis build be provisioned for request ID " + requestID)
}

/*
	TravisCICheckBuildStutus

	Travis API doc for build:
	https://developer.travis-ci.com/resource/build#Build
*/
func TravisCICheckBuildStutus(buildID string) (string, error) {
	fetchURL := BuildString(
		travisAPIBaseURL, "/build/", buildID,
	)

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

func TravisCIWaitTillBuildFinish(buildID string) (string, error) {
	// poll for up to 20 minutes
	MaxPollingCount := 12 * 20
	pollingCount := 0
	state := "received"

	for (state != "passed" && state != "failed") && pollingCount <= MaxPollingCount {
		pollingCount++
		time.Sleep(5 * time.Second)
		Logger("INFO", "Polling build ", buildID, " ...")

		state, checkStatusError := TravisCICheckBuildStutus(buildID)
		if checkStatusError != nil {
			Logger("WARN", "Got error while polling. State: ", state, "; error: ", checkStatusError.Error())
		} else {
			Logger("INFO", "Polled state: ", state)
		}
	}

	if state == "passed" || state == "failed" {
		return state, nil
	}

	return "", errors.New("Time out while polling TravisCI for build " + buildID)
}
