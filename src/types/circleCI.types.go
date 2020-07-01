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

// type CircleCIPipelineListResponseType struct {
// 	NextPageToken string                 `json:"next_page_token"`
// 	Items         []CircleCIPipelineType `json:"items[]"`
// }

type CircleCIWorkflowType struct {
	Status         string `json:"status"`
	PipelineID     string `json:"pipeline_id"`
	PipelineNumber int    `json:"pipeline_number"`
	Name           string `json:"name"`
	CreatedAt      string `json:"created_at"`
}

// equivalent to build in travis
type CircleCIPipelineWorkflowListResponseType struct {
	Workflows []CircleCIWorkflowType `json:"items"`
}
