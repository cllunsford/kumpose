package kumpose

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/project"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func CompToDeployment(proj *project.Project) ([]byte, error) {
	var kubeData []byte

	for _, name := range proj.ServiceConfigs.Keys() {
		s, _ := proj.ServiceConfigs.Get(name)

		con, _ := kubeContainer(name, s)

		dep := &extensions.Deployment{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
			},
			ObjectMeta: api.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"service": name},
			},
			Spec: extensions.DeploymentSpec{
				Replicas: 1,
				Template: api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{"service": name},
					},
					Spec: api.PodSpec{
						Containers: []api.Container{
							con,
						},
					},
				},
			},
		}

		data, err := json.MarshalIndent(dep, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal rc to json %v", err)
		} else {
			kubeData = append(kubeData, data...)
		}
	}

	return kubeData, nil
}

func CompToKube(proj *project.Project) ([]byte, error) {
	var kubeData []byte

	for _, name := range proj.ServiceConfigs.Keys() {
		s, _ := proj.ServiceConfigs.Get(name)

		con, _ := kubeContainer(name, s)

		rc := &api.ReplicationController{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "ReplicationController",
				APIVersion: "v1",
			},
			ObjectMeta: api.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"service": name},
			},
			Spec: api.ReplicationControllerSpec{
				Replicas: 1,
				Selector: map[string]string{"service": name},
				Template: &api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{"service": name},
					},
					Spec: api.PodSpec{
						Containers: []api.Container{
							con,
						},
					},
				},
			},
		}

		data, err := json.MarshalIndent(rc, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal rc to json %v", err)
		} else {
			kubeData = append(kubeData, data...)
		}
	}
	return kubeData, nil
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

	return container, nil
}

func kubePorts(ports []string) ([]api.ContainerPort, error) {
	var kubePorts []api.ContainerPort

	for _, port := range ports {
		kubePort := api.ContainerPort{}

		if strings.Contains(port, "/") {
			//Protocol is being specified
			protoParts := strings.Split(port, "/")
			port = protoParts[0]
			switch protoParts[1] {
			case "tcp":
				kubePort.Protocol = api.ProtocolTCP
			case "udp":
				kubePort.Protocol = api.ProtocolUDP
			default:
				log.Fatalf("Invalid protocol for container port %s", protoParts[1])
			}
		}

		if strings.Contains(port, ":") {
			//Host Port is being specified
			portParts := strings.Split(port, ":")
			port = portParts[1]
			if hostPort, err := strconv.ParseInt(portParts[0], 10, 32); err != nil {
				return kubePorts, err
			} else {
				kubePort.HostPort = int32(hostPort)
			}
		}
		if contPort, err := strconv.ParseInt(port, 10, 32); err != nil {
			return kubePorts, err
		} else {
			kubePort.ContainerPort = int32(contPort)
		}

		kubePorts = append(kubePorts, kubePort)
	}

	return kubePorts, nil
}

func kubeEnv(environment []string) ([]api.EnvVar, error) {
	var kubeEnv []api.EnvVar

	for _, e := range environment {
		if strings.Contains(e, "=") {
			parts := strings.Split(e, "=")
			kubeEnv = append(kubeEnv, api.EnvVar{
				Name:  parts[0],
				Value: parts[1],
			})
		} else {
			kubeEnv = append(kubeEnv, api.EnvVar{
				Name: e,
			})
		}
	}

	return kubeEnv, nil
}

func Run(composeFile string, target string) ([]byte, error) {

	proj := project.NewProject(
		&project.Context{
			ProjectName:  "k",
			ComposeFiles: []string{composeFile},
		},nil,nil)

	if err := proj.Parse(); err != nil {
		log.Fatal("Error parsing compose file")
		return []byte{}, err
	}

	switch target {
	case "rc":
		return CompToKube(proj)
	default:
		return CompToDeployment(proj)
	}
}
