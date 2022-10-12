package myErrors

import (
	"fmt"
	"net/http"
)

type authorizationError myError

func (a *authorizationError) Error() string {
	return (*myError)(a).Error()
}

func (a *authorizationError) GetStatusCode() int {
	return (*myError)(a).GetStatusCode()
}

func (a *authorizationError) GetEvent() string {
	return (*myError)(a).GetEvent()
}

func (a *authorizationError) GetReason() string {
	return (*myError)(a).GetReason()
}

func NewAuthorizationError(err error, event, reason string) *authorizationError {
	return (*authorizationError)(NewMyError(err, http.StatusUnauthorized, event, reason))
}

// the error will is fmt.Errorf("event: %s, reason: %s", event, reason)
func NewSimpleAuthorizationError(event, reason string) *authorizationError {
	return NewAuthorizationError(fmt.Errorf("event: %s, reason: %s", event, reason), event, reason)
}
