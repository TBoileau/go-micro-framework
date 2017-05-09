package framework

import (
	"net/http"
	"reflect"

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

		arguments := []reflect.Value{}
		for _, arg := range mux.Vars(request) {
			arguments = append(arguments, reflect.ValueOf(arg))
		}
		controller.MethodByName(route.Action).Call(arguments)
	})
}
