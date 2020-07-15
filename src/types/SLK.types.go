package types

type SLKS3JobResponseType struct {
	Progress int    `json:"progress"`
	Status   string `json:"status"`

	// if successfully provisioned new s3 job
	ID       string `json:"id"`
	Attempts int    `json:"attempts"`
	JobError string `json:"jobError"`

	// if failed to provision new s3 job
	Error string `json:"error"`
}
