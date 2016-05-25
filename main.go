package main

import (
	"flag"
	"fmt"

	lib "github.com/cllunsford/kumpose/kumpose"
)

var (
	composeFile = flag.String("compose-file", "docker-compose.yml", "Path to docker-compose file")
)

func init() {
	flag.StringVar(composeFile, "f", "docker-compose.yml", "Path to docker-compose file")
}

func main() {
	flag.Parse()

	if data, err := lib.Run(*composeFile, "deployment"); err == nil {
		fmt.Println(string(data))
	}
}
