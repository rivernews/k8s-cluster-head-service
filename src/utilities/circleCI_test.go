package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCircleCIBuildStatus(t *testing.T) {
	pipelineID := "5c9ab317-3f41-4851-a3de-e5fb119da8e6"
	pipelineWorkflows := FetchCircleCIBuildStatus(pipelineID)

	assert.Len(t, pipelineWorkflows.Workflows, 2, "Should have 2 workflows")

	latestWorkflow := pipelineWorkflows.Workflows[0]

	assert.Equal(t, latestWorkflow.Name, "build-master", "The name should be `build-master`")
	assert.Equal(t, latestWorkflow.PipelineNumber, 243, "The pipeline number should be 243")
	assert.Equal(t, latestWorkflow.Status, "canceled", "The status should be canceled")

}
