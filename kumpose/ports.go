package kumpose

import (
	"log"
	"strconv"
	"strings"

	"k8s.io/kubernetes/pkg/api"
)

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
