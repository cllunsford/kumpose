package kumpose

import (
	"encoding/json"
	"log"
	"os"
	"strings"

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
						Hostname: s.Hostname,
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

func kubeEnv(environment []string) ([]api.EnvVar, error) {
	var kubeEnv []api.EnvVar

	for _, e := range environment {
		envVar := api.EnvVar{}
		if strings.Contains(e, "=") {
			parts := strings.Split(e, "=")
			envVar.Name = parts[0]
			envVar.Value = parts[1]
		} else {
			envVar.Name = e
			v, ok := os.LookupEnv(e)
			if ok {
				envVar.Value = v
			}
		}
		kubeEnv = append(kubeEnv, envVar)
	}

	return kubeEnv, nil
}

func Run(composeFile string, target string) ([]byte, error) {

	proj := project.NewProject(
		&project.Context{
			ProjectName:  "k",
			ComposeFiles: []string{composeFile},
		}, nil, nil)

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
