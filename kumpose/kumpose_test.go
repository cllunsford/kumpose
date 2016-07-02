package kumpose

import (
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/stretchr/testify/assert"
)

func TestKubeEnv(t *testing.T) {
	assert := assert.New(t)

	env := []string{
		"TEST_VAR",
		"TEST_VAR_1=a value",
	}

	kubeEnv, err := kubeEnv(env)
	assert.NoError(err)
	assert.Equal(kubeEnv[0].Name, "TEST_VAR")
	assert.Equal(kubeEnv[0].Value, "")
	assert.Equal(kubeEnv[1].Name, "TEST_VAR_1")
	assert.Equal(kubeEnv[1].Value, "a value")
}
