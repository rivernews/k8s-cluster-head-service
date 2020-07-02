package types

type TravisCIBuildRequestResponseType struct {
	Builds []TravisCIBuild `json:"builds"`
}

type TravisCIBuild struct {
	ID        int    `json:"id"`
	State     string `json:"state"`
	StartedAt string `json:"started_at"`
}

type TravisCIRequestType struct {
	ID     int    `json:"id"`
	Branch string `json:"branch"`
}

// shape documetned at
// https://developer.travis-ci.com/resource/requests#create
type TravisCIRequestProvisionType struct {
	Type                   string              `json:"@type"`
	Request                TravisCIRequestType `json:"request"`
	RemainingRequestsCount int                 `json:"remaining_requests"`
}
