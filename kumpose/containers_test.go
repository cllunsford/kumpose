package kumpose

import (
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/stretchr/testify/assert"
)

func TestCompToContainer(t *testing.T) {
	assert := assert.New(t)

	sc := &config.ServiceConfig{
		Image: "nginx",
	}

	_, err := kubeContainer("name", sc)
	assert.NoError(err)
}
