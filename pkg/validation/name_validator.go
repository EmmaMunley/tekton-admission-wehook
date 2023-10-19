package validation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// pipelineNameValidator is a container for validating the name of pods
type pipelineNameValidator struct {
	Logger logrus.FieldLogger
}

type taskNameValidator struct {
	Logger logrus.FieldLogger
}

// pipelineNameValidator implements the pipelineValidator interface
var _ pipelineValidator = (*pipelineNameValidator)(nil)

// nameValidator implements the pipelineValidator interface
var _ taskValidator = (*taskNameValidator)(nil)

// Name returns the name of pipelineNameValidator
func (n pipelineNameValidator) Name() string {
	return "pipeline_name_validator"
}

// Name returns the name of taskNameValidator
func (n taskNameValidator) Name() string {
	return "task_name_validator"
}

// Validate inspects the name of a given pipeline and returns validation.
// The returned validation is only valid if the pipeline name does not contain some
// bad string.
func (n pipelineNameValidator) Validate(pipeline v1.Pipeline) (validation, error) {
	badString := "offensive"

	if strings.Contains(pipeline.Name, badString) {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("pipeline name contains %q", badString),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid name"}, nil
}

// Validate inspects the name of a given task and returns validation.
// The returned validation is only valid if the task name does not contain some
// bad string.
func (n taskNameValidator) Validate(task v1.Task) (validation, error) {
	badString := "offensive"

	if strings.Contains(task.Name, badString) {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("task name contains %q", badString),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid name"}, nil
}
