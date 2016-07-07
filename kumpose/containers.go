package kumpose

import (
	"fmt"
	"strings"

	"github.com/docker/libcompose/config"
	"k8s.io/kubernetes/pkg/api"
)

func kubePodSpec(name string, sc *config.ServiceConfig) (api.PodSpec, error) {
	ps := api.PodSpec{
		Hostname: sc.Hostname,
	}

	if sc.Restart != "" {
		if rp, err := kubeRestartPolicy(sc.Restart); err == nil {
			ps.RestartPolicy = rp
		}
	}

	return ps, nil
}

func kubeContainer(name string, sc *config.ServiceConfig) (api.Container, error) {
	container := api.Container{
		Name:            name,
		Image:           sc.Image,
		Command:         sc.Entrypoint,
		Args:            sc.Command,
		ImagePullPolicy: api.PullIfNotPresent,
		TTY:             sc.Tty,
		WorkingDir:      sc.WorkingDir,
	}
	if ports, err := kubePorts(sc.Ports); err == nil {
		container.Ports = ports
	}
	if env, err := kubeEnv(sc.Environment); err == nil {
		container.Env = env
	}
	if sec, err := kubeSecurityContext(sc); err == nil {
		container.SecurityContext = sec
	}

	return container, nil
}

func kubeSecurityContext(sc *config.ServiceConfig) (*api.SecurityContext, error) {
	securityContext := &api.SecurityContext{
		Capabilities: &api.Capabilities{
			Add:  stringsToCapability(sc.CapAdd),
			Drop: stringsToCapability(sc.CapDrop),
		},
		Privileged: &sc.Privileged,
	}
	return securityContext, nil
}

func kubeRestartPolicy(r string) (api.RestartPolicy, error) {
	switch strings.Split(r, ":")[0] {
	case "always":
		return api.RestartPolicyAlways, nil
	case "on-failure":
		return api.RestartPolicyOnFailure, nil
	case "no":
		return api.RestartPolicyNever, nil
	default:
		return api.RestartPolicyNever, fmt.Errorf("Restart policy not implemented: %s", r)
	}
}

func stringsToCapability(input []string) []api.Capability {
	caps := make([]api.Capability, len(input))
	for i, s := range input {
		caps[i] = api.Capability(s)
	}
	return caps
}
