package utilities

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTravisCICheckBuildStutus(t *testing.T) {
	buildID := "171515697"
	status, err := TravisCICheckBuildStutus(buildID)
	if err != nil {
		log.Print(err)
	}
	assert.Equal(t, status, "passed", "The status should be passed")
}
