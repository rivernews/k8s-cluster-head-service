package types

type CircleCIKubernetesClusterProjectPipelineParameters struct {
	DropletSize string `json:"kubernetes-cluster-droplet-size"`
}

type CircleCIRequestType struct {
	Branch     string                                             `json:"branch"`
	Parameters CircleCIKubernetesClusterProjectPipelineParameters `json:"parameters"`
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
