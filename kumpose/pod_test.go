package kumpose

import (
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestKubePodSpec(t *testing.T) {
	assert := assert.New(t)

	sc := &config.ServiceConfig{
		Image:   "nginx",
		Restart: "always",
	}

	ps, err := kubePodSpec("name", sc)
	assert.NoError(err)
	assert.Equal(api.RestartPolicyAlways, ps.RestartPolicy)
}

func TestKubePodSecurityContext(t *testing.T) {
	assert := assert.New(t)

	sc := &config.ServiceConfig{
		NetworkMode: "host",
	}

	psc, err := kubePodSecurityContext(sc)
	assert.NoError(err)
	assert.True(psc.HostNetwork)

	sc.NetworkMode = "bridge"
	psc, err = kubePodSecurityContext(sc)
	assert.NotNil(err)
	assert.False(psc.HostNetwork)
}

func TestKubeRestartPolicy(t *testing.T) {
	assert := assert.New(t)

	r, err := kubeRestartPolicy("always")
	assert.NoError(err)
	assert.Equal(api.RestartPolicyAlways, r)

	r, err = kubeRestartPolicy("on-failure")
	assert.NoError(err)
	assert.Equal(api.RestartPolicyOnFailure, r)

	r, err = kubeRestartPolicy("no")
	assert.NoError(err)
	assert.Equal(api.RestartPolicyNever, r)

	r, err = kubeRestartPolicy("foo")
	assert.NotNil(err)
	assert.Equal(api.RestartPolicyNever, r)
}
