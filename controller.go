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
}

func (controller Controller) Get() map[string]string {
	return mux.Vars(controller.Request)
}

func (controller Controller) Post() map[string]string {
	return mux.Vars(controller.Request)
}

func (controller Controller) Json(data interface{}) {
	json.NewEncoder(controller.ResponseWriter).Encode(data)
}

func (controller Controller) Print(format string, a ...interface{}) {
	fmt.Fprintf(controller.ResponseWriter, format, a...)
}
