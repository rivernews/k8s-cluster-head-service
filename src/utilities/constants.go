package utilities

import "os"

var RequestFromSlackTokenCredential, RequestFromSlackTokenCredentialExists = os.LookupEnv("REQUEST_FROM_SLACK_TOKEN")
var CircleCiToken, CircleCiTokenExists = os.LookupEnv("CIRCLECI_TOKEN")
var TravisCIToken, TravisCITokenExists = os.LookupEnv("TRAVIS_TOKEN")
