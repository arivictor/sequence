package sequence

import (
	"fmt"
	"os/exec"
)

type Job struct {
	Name        string   `yaml:"name"`
	Command     string   `yaml:"command"`
	ExitOnError bool     `yaml:"exit_on_error"`
	ErrorHook   string   `yaml:"error_hook"`
	SuccessHook string   `yaml:"success_hook"`
	Skip        bool     `yaml:"skip"`
	DependsOn   []string `yaml:"depends_on"`
}

func (job *Job) dependenciesSatisfied(executedJobs map[string]bool) bool {
	for _, dep := range job.getDependencies() {
		if success, exists := executedJobs[dep]; !exists || !success {
			return false
		}
	}
	return true
}

func (j Job) getDependencies() []string {
	return j.DependsOn
}

func (j Job) getSuccessHook() string {
	return j.SuccessHook
}

func (j Job) getErrorHook() string {
	return j.ErrorHook
}

func (j Job) getCommand() string {
	return j.Command
}

func (j Job) getName() string {
	return j.Name
}

func (j Job) getExitOnError() bool {
	return j.ExitOnError
}

func (j Job) getSkip() bool {
	return j.Skip
}

func (j Job) execute() error {
	cmd := exec.Command("bash", "-c", j.getCommand())
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		return err
	}

	return nil
}
