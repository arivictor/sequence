package sequence

import (
	"fmt"
	"os/exec"
)

type Hook struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func (h Hook) getName() string {
	return h.Name
}

func (h Hook) getCommand() string {
	return h.Command
}

func (h Hook) execute() error {
	cmd := exec.Command("bash", "-c", h.getCommand())
	o, err := cmd.CombinedOutput()
	fmt.Println(string(o))

	if err != nil {
		return err
	}
	return nil
}
