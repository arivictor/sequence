package sequence

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Jobs          []Job          `yaml:"jobs"`
	ErrorHandlers []ErrorHandler `yaml:"error_handlers"`
}

func (config Config) Execute() {
	executedJobs := make(map[string]bool) // Track executed jobs
	for _, job := range config.Jobs {
		if job.Skip {
			log.Printf("Skipping job '%s' as it's marked to be skipped.\n", job.Name)
			continue
		}

		// Check if all dependencies have been executed successfully
		allDependenciesSatisfied := true
		for _, dependency := range job.DependsOn {
			if !executedJobs[dependency] {
				allDependenciesSatisfied = false
				log.Printf("Skipping job '%s' because dependency '%s' has not been executed successfully.\n", job.Name, dependency)
				break
			}
		}

		if allDependenciesSatisfied {
			if executeJob(job, config) {
				executedJobs[job.Name] = true
			}
		}
	}
}

func NewConfig(configPath string) (Config, error) {
	c := Config{}
	y, err := os.ReadFile(configPath)
	if err != nil {
		return c, fmt.Errorf("error reading YAML file: %s", err)
	}

	err = yaml.Unmarshal(y, &c)
	if err != nil {
		return c, fmt.Errorf("error parsing YAML file: %s", err)
	}

	isValid := c.validate()
	if !isValid {
		return c, fmt.Errorf("invalid config file: %s", configPath)
	}
	return c, nil
}

func (config Config) validate() bool {
	return areJobNamesUnique(config.Jobs)
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
