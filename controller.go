package framework

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Server         *Server
}

func (controller Controller) GetMiddleware(middleware string) interface{} {
	return controller.Server.Middlewares[middleware].Interface()
}

func (controller Controller) Get() map[string]string {
	return mux.Vars(controller.Request)
}

func (controller Controller) Post() map[string]string {
	return mux.Vars(controller.Request)
}

func (controller Controller) Json(data interface{}) {
	controller.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(controller.ResponseWriter).Encode(data)
}

func (controller Controller) Print(format string, a ...interface{}) {
	fmt.Fprintf(controller.ResponseWriter, format, a...)
}
