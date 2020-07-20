package utilities

import (
	"errors"
	"strconv"
	"time"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"
)

func getSLKAPIBaseURL() string {
	if Debug {
		return "http://host.docker.internal:8080"
	}

	return "https://slack.api.shaungc.com"
}

func SLKCheckS3JobStatus() (types.SLKS3JobResponseType, error) {
	fetchURL := BuildString(
		getSLKAPIBaseURL(),
		"/queues/s3-orgs-job",
	)

	var s3JobResponse types.SLKS3JobResponseType
	_, _, fetchErr := Fetch(FetchOption{
		Method: "POST",
		URL:    fetchURL,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/json"},
		},
		PostData: map[string]string{
			"token":                   SLKSlackTokenOutgoingLaunch,
			"singleton":               "true",
			"keepAliveK8sHeadService": "true",
		},
		responseStore: &s3JobResponse,
	})
	// error handling
	if fetchErr != nil {
		return s3JobResponse, fetchErr
	}

	// return parse response
	return s3JobResponse, nil
}

func SLKWaitTillS3JobFinish() (types.SLKS3JobResponseType, error) {
	// polling for up to 1 day
	MaxPollingCount := 6 * 60 * 24
	pollingCount := 0

	for pollingCount <= MaxPollingCount {
		pollingCount++
		time.Sleep(10 * time.Second)
		Logger("INFO", "Polling s3 job...")

		s3JobMeta, checkStatusError := SLKCheckS3JobStatus()

		if checkStatusError != nil {
			Logger("WARN", "Got error while polling s3 job. Error: ", checkStatusError.Error())
		} else {
			Logger("INFO", "Polled progress=", strconv.Itoa(s3JobMeta.Progress), "; status=", s3JobMeta.Status, "; error=", s3JobMeta.Error, "; id=", s3JobMeta.ID, "; attempts=", strconv.Itoa(s3JobMeta.Attempts), "\n")
		}

		// terminate condition
		if s3JobMeta.Status == "failed" {
			return s3JobMeta, errors.New(s3JobMeta.JobError)
		} else if s3JobMeta.Status == "completed" {
			return s3JobMeta, nil
		}
	}

	return types.SLKS3JobResponseType{}, errors.New("Time out while polling s3 job to finish")
}
