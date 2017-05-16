package framework

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	yaml "gopkg.in/yaml.v2"
)

type Server struct {
	ProjectPath string
	Config      struct {
		StaticDic     string `yaml:"static_dir"`
		StaticPath    string `yaml:"static_path"`
		Port          string `yaml:"port"`
		RoutingPrefix string `yaml:"routing_prefix"`
		RoutingPath   string `yaml:"routing_path"`
		SourcesPath   string `yaml:"src_path"`
	} `yaml:"framework"`
	Parameters  interface{}
	Bundles     map[string]map[string]reflect.Value
	Middlewares map[string]reflect.Value
	Routing     map[string]Route
}

var server Server

func (server *Server) Initialize() {
	server.ProjectPath, _ = os.Getwd()
	server.Routing = make(map[string]Route)
	server.Bundles = make(map[string]map[string]reflect.Value)
	server.Middlewares = make(map[string]reflect.Value)
	server.LoadConfig()
	server.LoadRouting(Route{
		Resource: server.Config.RoutingPath,
		Prefix:   server.Config.RoutingPrefix,
	})
}

func CreateServer() *Server {
	server.Initialize()
	return &server
}

func (server *Server) RegisterBundles(bundles ...BundleInterface) {
	for _, bundle := range bundles {
		server.Bundles[bundle.GetName()] = bundle.Register()
	}
}

func (server *Server) RegisterMiddlewares(middlewares ...MiddlewareInterface) {
	for _, middleware := range middlewares {
		server.Middlewares[middleware.GetName()] = middleware.Register()
	}
}

func (server *Server) Run() {
	serv := negroni.Classic()
	mx := mux.NewRouter()
	for _, route := range server.Routing {
		route.Handle(mx)
	}
	if err := server.Config.StaticDic; err != "" {
		mx.PathPrefix(server.Config.StaticPath).Handler(http.FileServer(http.Dir(server.Config.StaticDic)))
	}
	serv.UseHandler(mx)
	serv.Run(":" + server.Config.Port)
}

func (server *Server) LoadConfig() {
	yamlFile, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.Panic("No routing file : config/config.yml")
	}

	err = yaml.Unmarshal(yamlFile, &server)
	if err != nil {
		log.Panic("Error on parsing routing file : config/config.yml")
	}
}

func (server *Server) LoadRouting(parent Route) {
	yamlFile, err := ioutil.ReadFile(parent.Resource)
	if err != nil {
		log.Panic("No routing file : " + parent.Resource)
	}

	m := make(map[string]map[string]string)
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Panic("Error on parsing routing file : " + parent.Resource)
	}

	for name, def := range m {
		if prefix, ok := def["prefix"]; ok {
			if bundle, ok := def["bundle"]; ok {
				if resource, ok := def["resource"]; ok {
					if prefix == "/" {
						prefix = ""
					}
					route := Route{
						Name:     name,
						Prefix:   parent.Prefix + prefix,
						Bundle:   bundle,
						Resource: resource,
						Parent:   &parent,
					}
					server.LoadRouting(route)
				}
			}
		}
		if path, ok := def["path"]; ok {
			if controller, ok := def["controller"]; ok {
				if action, ok := def["action"]; ok {
					if path == "/" {
						path = ""
					}
					route := Route{
						Name:       name,
						Action:     action,
						Path:       parent.Prefix + path,
						Controller: controller,
						Parent:     &parent,
					}
					server.Routing[parent.Prefix+path] = route
				}
			}
		}
	}
}
