package errors

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ErrorIs returns true if an error satisfies a particular condition.
type ErrorIs func(err error) bool

// IgnoreAny ignores errors that satisfy any of the supplied ErrorIs functions
// by returning nil. Errors that do not satisfy any of the supplied functions
// are returned unmodified.
func IgnoreAny(err error, is ...ErrorIs) error {
	for _, f := range is {
		if f(err) {
			return nil
		}
	}

	return err
}

// Ignore any errors that satisfy the supplied ErrorIs function by returning
// nil. Errors that do not satisfy the supplied function are returned unmodified.
func Ignore(err error, is ErrorIs) error {
	return IgnoreAny(err, is)
}

// IsNamespaceTerminating returns true if the namespace is terminating.
func IsNamespaceTerminating(err error) bool {
	return apierrors.HasStatusCause(err, corev1.NamespaceTerminatingCause)
}

type JobPodNotFoundError struct {
	JobName string
}

func (e *JobPodNotFoundError) Error() string {
	return fmt.Sprintf("no pods found for job %q", e.JobName)
}

func IsJobPodNotFound(err error) bool {
	var notFoundErr *JobPodNotFoundError
	return errors.As(err, &notFoundErr)
}

func NewJobPodNotFound(jobName string) *JobPodNotFoundError {
	return &JobPodNotFoundError{JobName: jobName}
}

type ScanJobContainerWaitingError struct {
	State corev1.ContainerStateWaiting
}

func (e *ScanJobContainerWaitingError) Error() string {
	return e.State.Message
}

func IsScanJobContainerWaiting(err error) bool {
	var waitingErr *ScanJobContainerWaitingError
	return errors.As(err, &waitingErr)
}

func NewScanJobContainerWaiting(state corev1.ContainerStateWaiting) *ScanJobContainerWaitingError {
	return &ScanJobContainerWaitingError{State: state}
}
