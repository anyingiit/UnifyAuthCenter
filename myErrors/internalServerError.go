package myErrors

import (
	"fmt"
	"net/http"
)

type internalError myError

func (i *internalError) Error() string {
	return (*myError)(i).Error()
}

func (i *internalError) GetStatusCode() int {
	return (*myError)(i).GetStatusCode()
}

func (i *internalError) GetEvent() string {
	return (*myError)(i).GetEvent()
}

func (i *internalError) GetReason() string {
	return (*myError)(i).GetReason()
}

func NewInternalError(err error, event, reason string) *internalError {
	return (*internalError)(NewMyError(err, http.StatusInternalServerError, event, reason))
}

// the error will is fmt.Errorf("event: %s, reason: %s", event, reason)
func NewSimpleInternalError(event, reason string) *internalError {
	return NewInternalError(fmt.Errorf("event: %s, reason: %s", event, reason), event, reason)
}
