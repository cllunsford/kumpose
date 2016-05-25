package kumpose

import (
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestCompToContainer(t *testing.T) {
	assert := assert.New(t)

	sc := &config.ServiceConfig{
		Image: "nginx",
	}

	_, err := kubeContainer("name", sc)
	assert.NoError(err)
}

func TestKubePortProtocols(t *testing.T) {
	assert := assert.New(t)

	ports := []string{
		"80",
		"81/tcp",
		"82/udp",
	}

	kubePorts, err := kubePorts(ports)
	assert.NoError(err)
	assert.Equal(kubePorts[1].Protocol, api.ProtocolTCP)
	assert.Equal(kubePorts[2].Protocol, api.ProtocolUDP)
}

func TestKubePortHostPorts(t *testing.T) {
	assert := assert.New(t)

	ports := []string{
		"80:8080",
		"443:8081/tcp",
		"444:8082/udp",
	}

	kubePorts, err := kubePorts(ports)
	assert.NoError(err)
	assert.Equal(kubePorts[0].ContainerPort, int32(8080))
	assert.Equal(kubePorts[0].HostPort, int32(80))
	assert.Equal(kubePorts[1].ContainerPort, int32(8081))
	assert.Equal(kubePorts[1].HostPort, int32(443))
	assert.Equal(kubePorts[1].Protocol, api.ProtocolTCP)
	assert.Equal(kubePorts[2].ContainerPort, int32(8082))
	assert.Equal(kubePorts[2].HostPort, int32(444))
	assert.Equal(kubePorts[2].Protocol, api.ProtocolUDP)
}

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
