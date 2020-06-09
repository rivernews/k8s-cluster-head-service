package types

type CircleCIRequestType struct {
	Branch  string            `json:"branch"`
	EnvVars map[string]string `json:"parameters"`
}

type CircleCIPipelineType struct {
	ID        string `json:"id"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
}

type CircleCIPipelineListResponseType struct {
	NextPageToken string                 `json:"next_page_token"`
	Items         []CircleCIPipelineType `json:"items[]"`
}
