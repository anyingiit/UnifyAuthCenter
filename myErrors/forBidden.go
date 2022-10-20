package myErrors

import (
	"fmt"
	"net/http"
)

type forBiddenError myError

func (a *forBiddenError) Error() string {
	return (*myError)(a).Error()
}

func (a *forBiddenError) GetStatusCode() int {
	return (*myError)(a).GetStatusCode()
}

func (a *forBiddenError) GetEvent() string {
	return (*myError)(a).GetEvent()
}

func (a *forBiddenError) GetReason() string {
	return (*myError)(a).GetReason()
}

func NewForBiddenErrorError(err error, event, reason string) *forBiddenError {
	return (*forBiddenError)(NewMyError(err, http.StatusForbidden, event, reason))
}

// the error will is fmt.Errorf("event: %s, reason: %s", event, reason)
func NewSimpleForBiddenErrorError(event, reason string) *forBiddenError {
	return NewForBiddenErrorError(fmt.Errorf("event: %s, reason: %s", event, reason), event, reason)
}
