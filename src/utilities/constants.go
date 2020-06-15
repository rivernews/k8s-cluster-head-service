package utilities

import (
	"os"
	"strings"
)

var RequestFromSlackTokenCredential, RequestFromSlackTokenCredentialExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")
var SendSlackURL, SendSlackURLExists = os.LookupEnv("SLACK_INCOMING_URL")
var CircleCiToken, CircleCiTokenExists = os.LookupEnv("CIRCLECI_TOKEN")
var TravisCIToken, TravisCITokenExists = os.LookupEnv("TRAVIS_TOKEN")
var Debug = getBoolEnvVarHelper("DEBUG")

var LargeDroplet = "s-4vcpu-8gb"
var MediumDroplet = "s-2vcpu-4gb"
var SmallDroplet = "s-1vcpu-3gb"

// getEnvVarHelper - don't care about no value when getting env var.
// Do not use this for credential, because we should always make sure credentials are available
// to avoid comparing to empty string when auth
func getEnvVarHelper(key string) string {
	return getEnvVarOrDefault(key, "")
}
func getBoolEnvVarHelper(key string) bool {
	value := strings.TrimSpace(strings.ToLower(getEnvVarHelper(key)))
	if value == "true" || value == "yes" || value == "1" {
		return true
	}
	return false
}

// getEnvVarOrDefault - must give a default value
func getEnvVarOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}