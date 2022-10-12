package myErrors

import (
	"fmt"
	"net/http"
)

type badRequestError myError

func (b *badRequestError) Error() string {
	return (*myError)(b).Error()
}

func (b *badRequestError) GetStatusCode() int {
	return (*myError)(b).GetStatusCode()
}

func (b *badRequestError) GetEvent() string {
	return (*myError)(b).GetEvent()
}

func (b *badRequestError) GetReason() string {
	return (*myError)(b).GetReason()
}

func NewBadRequestError(err error, event, reason string) *badRequestError {
	return (*badRequestError)(NewMyError(err, http.StatusBadRequest, event, reason))
}

// the error will is fmt.Errorf("event: %s, reason: %s", event, reason)
func NewSimpleBadRequestError(event, reason string) *badRequestError {
	return NewBadRequestError(fmt.Errorf("event: %s, reason: %s", event, reason), event, reason)
}
