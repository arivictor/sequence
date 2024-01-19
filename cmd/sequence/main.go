package main

import (
	"flag"
	"log"

	"github.com/arivictor/sequence"
)

func main() {
	var yamlPath string
	flag.StringVar(
		&yamlPath,
		"config",
		"config.yaml",
		"path to the YAML config file",
	)
	flag.Parse()

	c, err := sequence.NewConfig(yamlPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	c.Execute()
}
