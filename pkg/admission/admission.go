// Package admission handles kubernetes admissions,
// it takes admission requests and returns admission reviews
// including validating tekton pipelines and tasks
package admission

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tekton-webhook-admission/pkg/validation"

	"github.com/sirupsen/logrus"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

// Admitter is a container for admission business
type Admitter struct {
	Logger  *logrus.Entry
	Request *admissionv1.AdmissionRequest
}

// ValidatePipelineReview takes an admission request and validates the pipeline within
// it returns an admission review
func (a Admitter) ValidatePipelineReview() (*admissionv1.AdmissionReview, error) {
	pipeline, err := a.Pipeline()
	if err != nil {
		e := fmt.Sprintf("could not parse pipeline in admission review request: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	v := validation.NewValidator(a.Logger)
	val, err := v.ValidatePipeline(pipeline)
	if err != nil {
		e := fmt.Sprintf("could not validate pipeline: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	if !val.Valid {
		return reviewResponse(a.Request.UID, false, http.StatusForbidden, val.Reason), nil
	}

	return reviewResponse(a.Request.UID, true, http.StatusAccepted, "valid pipeline"), nil
}

// ValidateTaskReview takes an admission request and validates the task within
// it returns an admission review
func (a Admitter) ValidateTaskReview() (*admissionv1.AdmissionReview, error) {
	task, err := a.Task()
	if err != nil {
		e := fmt.Sprintf("could not parse task in admission review request: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	v := validation.NewValidator(a.Logger)
	val, err := v.ValidateTask(task)
	if err != nil {
		e := fmt.Sprintf("could not validate task: %v", err)
		return reviewResponse(a.Request.UID, false, http.StatusBadRequest, e), err
	}

	if !val.Valid {
		return reviewResponse(a.Request.UID, false, http.StatusForbidden, val.Reason), nil
	}

	return reviewResponse(a.Request.UID, true, http.StatusAccepted, "valid task"), nil
}

// Pipeline extracts a pipeline from an admission request
func (a Admitter) Pipeline() (v1.Pipeline, error) {
	if a.Request.Kind.Kind != "Pipeline" {
		return v1.Pipeline{}, fmt.Errorf("only pipelines are supported here")
	}
	p := v1.Pipeline{}
	if err := json.Unmarshal(a.Request.Object.Raw, &p); err != nil {
		return v1.Pipeline{}, err
	}
	return p, nil
}

// Task extracts a task from an admission request
func (a Admitter) Task() (v1.Task, error) {
	if a.Request.Kind.Kind != "Task" {
		return v1.Task{}, fmt.Errorf("only tasks are supported here")
	}
	t := v1.Task{}
	if err := json.Unmarshal(a.Request.Object.Raw, &t); err != nil {
		return v1.Task{}, err
	}
	return t, nil
}

func reviewResponse(uid types.UID, allowed bool, httpCode int32,
	reason string) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

// patchReviewResponse builds an admission review with given json patch
func patchReviewResponse(uid types.UID, patch []byte) (*admissionv1.AdmissionReview, error) {
	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}, nil
}
