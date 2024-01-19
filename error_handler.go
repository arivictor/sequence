package sequence

import (
	"fmt"
	"log"
	"os/exec"
)

type ErrorHandler struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func handleError(job Job, config Config) {
	for _, errorHandler := range config.ErrorHandlers {
		if errorHandler.Name == job.ErrorHandler {
			fmt.Printf("Executing Error Handler: %s for Job: %s\n", errorHandler.Name, job.Name)
			handlerCmd := exec.Command("bash", "-c", errorHandler.Command)
			handlerOutput, handlerErr := handlerCmd.CombinedOutput()
			fmt.Println(string(handlerOutput))

			if handlerErr != nil {
				log.Printf("Error executing error handler '%s': %s\n", errorHandler.Name, handlerErr)
			}
			break
		}
	}
}
