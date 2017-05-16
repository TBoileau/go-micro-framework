package framework

import (
	"io/ioutil"
	"log"
	"reflect"
)

type Middleware struct {
	ConfigFile string
}

type MiddlewareInterface interface {
	GetName() string
	Register() reflect.Value
}

func (middleware *Middleware) LoadParameters() (yamlFile []byte) {
	yamlFile, err := ioutil.ReadFile(middleware.ConfigFile)
	if err != nil {
		log.Panic("No config file : " + middleware.ConfigFile)
	}
	return
}
