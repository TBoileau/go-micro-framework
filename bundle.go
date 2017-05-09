package framework

import "reflect"

type Bundle struct {
}

type BundleInterface interface {
	GetName() string
	Register() map[string]reflect.Value
}
