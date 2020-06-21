package utilities

import (
	"net/url"
	"strings"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"

	"github.com/gin-gonic/gin"
)

var travisAPIBaseURL = "https://api.travis-ci.com"

/*
	TravisCITriggerSLKHelper - ede
*/
func TravisCITriggerSLKHelper(c *gin.Context, parsedSlackRequest types.SlackRequestType) {
	encodedProjectSlug := url.QueryEscape("rivernews/slack-middleware-server")

	// build url
	var urlBuilder strings.Builder
	urlBuilder.WriteString(travisAPIBaseURL)
	// endpoint
	urlBuilder.WriteString("/repo/")
	urlBuilder.WriteString(encodedProjectSlug)
	urlBuilder.WriteString("/requests")

	_, fetchedMessage := Fetch(FetchOption{
		Method: "POST",
		URL:    urlBuilder.String(),
		Headers: map[string][]string{
			"Content-Type":       {"application/json"},
			"Accept":             {"application/json"},
			"Travis-API-Version": {"3"},
			"Authorization":      {"token " + TravisCIToken},
		},
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
