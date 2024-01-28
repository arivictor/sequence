package main

import (
	"flag"

	"github.com/arivictor/sequence"
)

func main() {
	var filepath string
	flag.StringVar(
		&filepath,
		"workflow",
		"workflow.yaml",
		"path to the YAML workflow file",
	)
	flag.Parse()

	w, err := sequence.NewWorkflow(filepath)
	if err != nil {
		panic(err)
	}

	err = w.Execute()
	if err != nil {
		panic(err)
	}
}
