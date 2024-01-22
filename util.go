package sequence

import (
	"os"

	"gopkg.in/yaml.v2"
)

func readFile(filepath string) ([]byte, error) {
	d, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func parseYaml(data []byte) (*Workflow, error) {
	w := Workflow{}
	err := yaml.Unmarshal(data, &w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
