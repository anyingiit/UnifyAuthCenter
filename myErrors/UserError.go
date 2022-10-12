package myErrors

type UserError interface {
	error
	GetStatusCode() int
	GetEvent() string
	GetReason() string
}
