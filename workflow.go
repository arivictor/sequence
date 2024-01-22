package sequence

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Workflow struct {
	Jobs    []Job  `yaml:"jobs"`
	Hooks   []Hook `yaml:"hooks"`
	hookMap map[string]*Hook
}

func (w *Workflow) getJobs() []Job {
	return w.Jobs
}

func (w *Workflow) getHooks() []Hook {
	return w.Hooks
}

func NewWorkflow(filepath string) (*Workflow, error) {
	d, err := readFile(filepath)
	if err != nil {
		return nil, err
	}

	w, err := parseYaml(d)
	if err != nil {
		return nil, err
	}

	w.buildHookMap()

	err = w.validate()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Workflow) Execute() error {
	executedJobs := make(map[string]bool) // Track executed jobs
	for _, job := range w.getJobs() {

		if job.getSkip() {
			logrus.WithFields(logrus.Fields{
				"type":   "job",
				"name":   job.getName(),
				"action": "skip",
			}).Info("property skip is set to true")
			continue
		}

		if !job.dependenciesSatisfied(executedJobs) {
			logrus.WithFields(logrus.Fields{
				"type":   "job",
				"name":   job.getName(),
				"action": "skip",
				"deps":   job.getDependencies(),
			}).Info("job dependencies not satisfied")
			continue
		}

		err := w.executeJob(&job)
		if err != nil {
			executedJobs[job.Name] = false
			if job.getExitOnError() {
				logrus.WithFields(logrus.Fields{
					"type":   "job",
					"name":   job.getName(),
					"action": "exit",
				}).Info("property exit_on_error is set to true")

				return fmt.Errorf("exiting due to error in job '%s': %v", job.getName(), err)
			}
			continue
		}

		executedJobs[job.Name] = true
	}

	return nil
}

func (w *Workflow) executeJob(job *Job) error {
	logrus.WithFields(logrus.Fields{
		"type":   "job",
		"name":   job.getName(),
		"action": "execute",
	}).Info("executing command")

	err := job.execute()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"type":   "job",
			"name":   job.getName(),
			"action": "error",
			"error":  err,
		}).Error("job execution failed")

		if hook, exists := w.hookMap[job.getErrorHook()]; exists {
			logrus.WithFields(logrus.Fields{
				"type":   "hook",
				"name":   hook.getName(),
				"job":    job.getName(),
				"action": "execute",
			}).Info("executing hook: on_error")
			hook.execute()
		}

		return err

	} else if hook, exists := w.hookMap[job.getSuccessHook()]; exists {
		logrus.WithFields(logrus.Fields{
			"type":   "hook",
			"name":   hook.getName(),
			"job":    job.getName(),
			"action": "execute",
		}).Info("executing hook: on_success")
		hook.execute()
	}

	return nil
}

func (w *Workflow) buildHookMap() {
	w.hookMap = make(map[string]*Hook)
	for i := range w.getHooks() {
		hook := &w.Hooks[i]
		w.hookMap[hook.Name] = hook
	}
}

func (w *Workflow) validate() error {
	if !w.validateJobUniqueNames() {
		return fmt.Errorf("workflow contains duplicate job names")
	}

	if !w.validateJobHooksExists() {
		return fmt.Errorf("workflow contains job hooks that do not exist")
	}

	return nil
}

func (w *Workflow) validateJobUniqueNames() bool {
	names := make(map[string]struct{})
	for _, job := range w.Jobs {
		if _, exists := names[job.Name]; exists {
			return false
		}
		names[job.Name] = struct{}{}
	}

	return true
}

func (w *Workflow) validateJobHooksExists() bool {
	for _, job := range w.Jobs {
		if job.getErrorHook() != "" && !w.hookExists(job.getErrorHook()) {
			logrus.WithField("job", job.Name).Error("Job contains on_error hook that does not exist")
			return false
		}

		if job.getSuccessHook() != "" && !w.hookExists(job.getSuccessHook()) {
			logrus.WithField("job", job.Name).Error("Job contains on_success hook that does not exist")
			return false
		}
	}

	return true
}

func (w *Workflow) hookExists(hookName string) bool {
	_, exists := w.hookMap[hookName]

	return exists
}
