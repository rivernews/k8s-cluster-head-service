package utilities

import (
	"os"
	"strings"
)

var RequestFromSlackTokenCredential, RequestFromSlackTokenCredentialExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")
var SendSlackURL, SendSlackURLExists = os.LookupEnv("SLACK_INCOMING_URL")
var CircleCiToken, CircleCiTokenExists = os.LookupEnv("CIRCLECI_TOKEN")
var TravisCIToken, TravisCITokenExists = os.LookupEnv("TRAVIS_TOKEN")

// used for authenticating with SLK
var SLKSlackTokenOutgoingLaunch, SLKSlackTokenOutgoingLaunchExists = os.LookupEnv("SLK_SLACK_TOKEN_OUTGOING_LAUNCH")

var Debug = getBoolEnvVarHelper("DEBUG")

var LargeDroplet = "s-4vcpu-8gb"
var MediumDroplet = "s-2vcpu-4gb"
var SmallDroplet = "s-1vcpu-3gb"

func GetRedisURL() (string, bool) {
	var redisURL string
	redisURLExists := false
	if Debug {
		// provides 20 connections, 25MB and ? db
		// https://elements.heroku.com/addons/heroku-redis
		redisURL, redisURLExists = os.LookupEnv("REDIS_URL")
	} else {
		// provides 30 connections, 30MB and 1 db
		// https://elements.heroku.com/addons/rediscloud
		redisURL, redisURLExists = os.LookupEnv("REDISCLOUD_URL")
	}

	return redisURL, redisURLExists
}

func GetRedisDB() int {
	return 0
}

var LogLevelTypes = map[string]int{
	"DEBUG": 4,
	"INFO":  3,
	"WARN":  2,
	"ERROR": 1,
}

func GetLogLevel() (int, string) {
	if Debug {
		return LogLevelTypes["DEBUG"], "DEBUG"
	}

	var logLevel = getEnvVarOrDefault("LOG_LEVEL", "INFO")
	if value, exist := LogLevelTypes[logLevel]; exist {
		return value, logLevel
	}

	return LogLevelTypes["INFO"], "INFO"
}

func GetLogLevelValue() int {
	value, _ := GetLogLevel()
	return value
}

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
