package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Job struct {
	Name        string   `yaml:"name"`
	Command     string   `yaml:"command"`
	ExitOnError bool     `yaml:"exit_on_error"`
	OnError     string   `yaml:"on_error"`
	OnSuccess   string   `yaml:"on_success"`
	Skip        bool     `yaml:"skip"`
	DependsOn   []string `yaml:"depends_on"`
}

type ErrorHandler struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Config struct {
	Jobs            []Job          `yaml:"jobs"`
	ErrorHandlers   []ErrorHandler `yaml:"error_handlers"`
	SuccessHandlers []ErrorHandler `yaml:"success_handlers"`
}

func executeJob(job Job, config Config) error {
	cmd := exec.Command("bash", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		return err
	}

	return nil
}

func areJobNamesUnique(jobs []Job) bool {
	names := make(map[string]bool)
	for _, job := range jobs {
		if _, exists := names[job.Name]; exists {
			return false
		}
		names[job.Name] = true
	}

	return true
}

func handleJobError(job Job, config Config) error {
	for _, errorHandler := range config.ErrorHandlers {
		if errorHandler.Name == job.OnError {
			log.Printf("execute : '%s' : error_handler '%s'\n", job.Name, errorHandler.Name)
			handlerCmd := exec.Command("bash", "-c", errorHandler.Command)
			handlerOutput, handlerErr := handlerCmd.CombinedOutput()
			fmt.Println(string(handlerOutput))

			if handlerErr != nil {
				return handlerErr
			}
			break
		}
	}

	return nil
}

func handleJobSuccess(job Job, config Config) error {
	for _, successHandler := range config.SuccessHandlers {
		if successHandler.Name == job.OnSuccess {
			log.Printf("execute : '%s' : success_handler '%s'\n", job.Name, successHandler.Name)
			handlerCmd := exec.Command("bash", "-c", successHandler.Command)
			handlerOutput, handlerErr := handlerCmd.CombinedOutput()
			fmt.Println(string(handlerOutput))

			if handlerErr != nil {
				return handlerErr
			}
			break
		}
	}

	return nil
}

func main() {
	// Command line flags
	var configPath string
	flag.StringVar(
		&configPath,
		"config",
		"config.yaml",
		"path to the YAML config file",
	)
	flag.Parse()

	// Load config
	c := Config{}
	y, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error : reading YAML file : %s", err)
	}

	err = yaml.Unmarshal(y, &c)
	if err != nil {
		log.Fatalf("error : parsing YAML file : %s", err)
	}

	// Validation
	jobNamesUnique := areJobNamesUnique(c.Jobs)
	if !jobNamesUnique {
		log.Fatalf("error : job names are not unique : %s", configPath)
	}

	// Execute jobs
	executedJobs := make(map[string]bool) // Track executed jobs
	for _, job := range c.Jobs {
		if job.Skip {
			log.Printf("skip : '%s' : skip is true\n", job.Name)
			continue
		}

		// Check if all dependencies have been executed successfully
		allDependenciesSatisfied := true
		for _, dependency := range job.DependsOn {
			if !executedJobs[dependency] {
				allDependenciesSatisfied = false
				log.Printf("pass : '%s' : depends_on '%s'\n", job.Name, dependency)
				break
			}
		}

		if allDependenciesSatisfied {
			log.Printf("execute : '%s'\n", job.Name)
			jobSuccess := true
			err := executeJob(job, c)
			if err != nil {
				handleJobErr := handleJobError(job, c)
				if handleJobErr != nil {
					log.Printf("error : '%s' : %s\n", job.Name, handleJobErr)
				}

				jobSuccess = false
				log.Printf("error : '%s' : %s\n", job.Name, err)
				if job.ExitOnError {
					log.Printf("exit : '%s' : exit_on_error is true\n", job.Name)
					os.Exit(1)
				}
			} else {
				err = handleJobSuccess(job, c)
				if err != nil {
					log.Printf("error : '%s' : on_success : %s\n", job.Name, err)
				}
			}
			executedJobs[job.Name] = jobSuccess
		}
	}
}
