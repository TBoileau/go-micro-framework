package framework

import (
	"net/http"
	"reflect"
	"regexp"

	"github.com/gorilla/mux"
)

type Route struct {
	*mux.Route
	Name       string
	Prefix     string
	Resource   string
	Bundle     string
	Controller string
	Action     string
	Path       string
	Parent     *Route
}

func (r *Route) Handle(router *mux.Router) {
	router.HandleFunc(r.Path, func(response http.ResponseWriter, request *http.Request) {
		path, _ := mux.CurrentRoute(request).GetPathTemplate()
		route := server.Routing[path]
		controller := server.Bundles[route.Parent.Bundle][route.Controller]
		controller.Elem().FieldByName("ResponseWriter").Set(reflect.ValueOf(response))
		controller.Elem().FieldByName("Request").Set(reflect.ValueOf(request))
		controller.Elem().FieldByName("Server").Set(reflect.ValueOf(&server))
		re := regexp.MustCompile(`{(\w*)}*`)
		i := 0
		method := controller.Elem().MethodByName(route.Action)
		in := make([]reflect.Value, method.Type().NumIn())
		vars := mux.Vars(request)
		for _, group := range re.FindAllStringSubmatch(path, -1) {
			in[i] = reflect.ValueOf(vars[group[1]])
			i++
		}
		method.Call(in)
	})
}
