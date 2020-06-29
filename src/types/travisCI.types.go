package types

type TravisCIBuildRequestResponseType struct {
	Builds []TravisCIBuild `json:"builds"`
}

type TravisCIBuild struct {
	ID        int    `json:"id"`
	State     string `json:"state"`
	StartedAt string `json:"started_at"`
}
