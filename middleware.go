package framework

import "reflect"

type Middleware struct {
}

type MiddlewareInterface interface {
	GetName() string
	Register() reflect.Value
}
