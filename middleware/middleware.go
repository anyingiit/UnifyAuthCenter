package middleware

type MiddlewareFunc func() (originHandle appHandler, middleware appHandler)
