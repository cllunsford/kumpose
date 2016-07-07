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

	if psc, err := kubePodSecurityContext(sc); err == nil {
		ps.SecurityContext = &psc
	}

	return ps, nil
}

func kubePodSecurityContext(sc *config.ServiceConfig) (api.PodSecurityContext, error) {
	ctx := api.PodSecurityContext{}

	if sc.NetworkMode != "" {
		switch sc.NetworkMode {
		case "host":
			ctx.HostNetwork = true
		default:
			return ctx, fmt.Errorf("NetworkMode not implemented: %s", sc.NetworkMode)
		}
	}

	return ctx, nil
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
