package sequence

import (
	"fmt"
	"log"
	"os/exec"
)

type Job struct {
	Name         string   `yaml:"name"`
	Command      string   `yaml:"command"`
	ExitOnError  bool     `yaml:"exit_on_error"`
	ErrorHandler string   `yaml:"error_handler"`
	Skip         bool     `yaml:"skip"`
	DependsOn    []string `yaml:"depends_on"`
}

func executeJob(job Job, config Config) bool {
	fmt.Printf("Executing Job: %s\n", job.Name)

	// Execute the command
	cmd := exec.Command("bash", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		log.Printf("Error executing job '%s': %s\n", job.Name, err)
		handleError(job, config)

		if job.ExitOnError {
			log.Fatalf("Job '%s' failed, exit_on_error is true. Exiting...\n", job.Name)
		}
		return false
	}
	return true
}
