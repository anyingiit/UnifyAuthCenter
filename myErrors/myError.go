package myErrors

type myError struct {
	Err        error
	StatusCode int
	Event      string
	Reason     string
}

func (e *myError) Error() string {
	return e.Err.Error()
}

func (e *myError) GetStatusCode() int {
	return e.StatusCode
}

func (e *myError) GetEvent() string {
	return e.Event
}

func (e *myError) GetReason() string {
	return e.Reason
}

func NewMyError(err error, statusCode int, event, reason string) *myError {
	return &myError{
		Err:        err,
		StatusCode: statusCode,
		Event:      event,
		Reason:     reason,
	}
}
