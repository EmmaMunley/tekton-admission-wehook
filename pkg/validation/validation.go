package validation

import (
	"github.com/sirupsen/logrus"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/webhook/resourcesemantics"
)

var types = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
	// v1
	v1.SchemeGroupVersion.WithKind("Task"):     &v1.Task{},
	v1.SchemeGroupVersion.WithKind("Pipeline"): &v1.Pipeline{},
}

// Validator is a container for mutation
type Validator struct {
	Logger *logrus.Entry
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

// pipelineValidators is an interface used to group functions mutating pods
type pipelineValidator interface {
	Validate(v1.Pipeline) (validation, error)
	Name() string
}

type taskValidator interface {
	Validate(v1.Task) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePipeline returns true if a pipeline is valid
func (v *Validator) ValidatePipeline(pipeline v1.Pipeline) (validation, error) {
	var pipelineName string
	if pipeline.Name != "" {
		pipelineName = pipeline.Name
	} else {
		if pipeline.ObjectMeta.GenerateName != "" {
			pipelineName = pipeline.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pipeline_name", pipelineName)
	log.Print("delete me")

	// list of all validations to be applied to the pipeline
	validations := []pipelineValidator{
		pipelineNameValidator{v.Logger},
	}

	// apply all validations
	for _, v := range validations {
		var err error
		vp, err := v.Validate(pipeline)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid pipeline"}, nil
}

// ValidateTask returns true if a task is valid
func (v *Validator) ValidateTask(task v1.Task) (validation, error) {
	var taskName string
	if task.Name != "" {
		taskName = task.Name
	} else {
		if task.ObjectMeta.GenerateName != "" {
			taskName = task.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("task_name", taskName)
	log.Print("delete me")

	// list of all validations to be applied to the task
	validations := []taskValidator{
		taskNameValidator{v.Logger},
	}
	// apply all validations
	for _, v := range validations {
		var err error
		vp, err := v.Validate(task)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid task"}, nil
}
