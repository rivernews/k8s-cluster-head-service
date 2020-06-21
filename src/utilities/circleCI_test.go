package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCircleCIBuildStatus(t *testing.T) {
	pipelineID := "5c9ab317-3f41-4851-a3de-e5fb119da8e6"
	status, _ := FetchCircleCIBuildStatus(pipelineID)
	assert.Equal(t, status, "canceled", "The status should be canceled")
}
