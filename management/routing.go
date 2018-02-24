package management

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/proxy"
	"github.com/husobee/vestigo"
)

//APIRouter ...
type APIRouter struct {
	user           *User
	data           *data.DataStore
	proxy          *proxy.Proxy
	router         *vestigo.Router
	serviceManager *ServiceManager
}

//NewAPIRouter ...
func NewAPIRouter(store *data.DataStore, proxy *proxy.Proxy, services *ServiceManager) *APIRouter {
	user := NewUserService(store)

	return &APIRouter{
		user:           user,
		data:           store,
		proxy:          proxy,
		serviceManager: services,
	}
}

//StartManagementAPIListener ...
func (api *APIRouter) StartManagementAPIListener() {
	api.router = vestigo.NewRouter()
	api.router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Cache-Control", "X-Requested-With", "Content-Type"},
		AllowOrigin:  []string{"https://svcman.m.rdro.us", "http://192.168.1.12:4200"},
	})
	vestigo.CustomNotFoundHandlerFunc(api.NotFound)
	api.SetupRoutes()
	go func() {
		if err := http.ListenAndServe("localhost:1000", api.router); err != nil {
			common.CreateFailureResponse(err, "StartManagementAPIListener", 500)
		}
	}()
}

//NotFound ...
func (api *APIRouter) NotFound(resp http.ResponseWriter, r *http.Request) {
	common.WriteFailureResponse(errors.New("route "+r.URL.String()+" not found."), resp, "NotFound", 404)
}

//SetupRoutes ...
func (api *APIRouter) SetupRoutes() {
	api.router.Handle("/ws/log", common.RequestWrapper(common.Nothing, "GET", api.ws))
	api.router.Handle("/api/frost/wsauth", common.RequestWrapper(api.user.IsRoot, "POST", api.wsauth))

	api.router.Handle("/api/frost/user/login", common.RequestWrapper(common.Nothing, "POST", api.login))
	api.router.Handle("/api/frost/user/register", common.RequestWrapper(common.Nothing, "POST", api.register))

	api.router.Handle("/api/frost/service/get", common.RequestWrapper(api.user.IsRoot, "GET", api.getService))
	api.router.Handle("/api/frost/service/new", common.RequestWrapper(api.user.IsRoot, "POST", api.newService))
	api.router.Handle("/api/frost/service/delete", common.RequestWrapper(api.user.IsRoot, "DELETE", api.deleteService))
	api.router.Handle("/api/frost/service/update", common.RequestWrapper(api.user.IsRoot, "POST", api.updateService))
	api.router.Handle("/api/frost/services", common.RequestWrapper(api.user.IsRoot, "GET", api.services))

	api.router.Handle("/api/frost/process", common.RequestWrapper(api.user.IsRoot, "GET", api.process))
}
func (api *APIRouter) login(resp http.ResponseWriter, r *http.Request) {
	var response common.APIResponse
	response = api.user.ValidateLoginRequest(r)
	common.WriteAPIResponseStruct(resp, response)
}
func (api *APIRouter) register(resp http.ResponseWriter, r *http.Request) {
	var request AuthRequest
	body, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, &request); err == nil {
		common.WriteAPIResponseStruct(resp, api.user.NewUser(request.Username, request.Password))
	} else {
		common.WriteFailureResponse(fmt.Errorf("failed deserializing request body %s", err), resp, "register", 500)
	}
}
func (api *APIRouter) newService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, api.serviceManager.NewService(r))
}
func (api *APIRouter) deleteService(resp http.ResponseWriter, r *http.Request) {
	var name = r.URL.Query().Get("name")
	api.serviceManager.StopManagedService(name)
	common.WriteAPIResponseStruct(resp, api.serviceManager.DeleteService(name))
}
func (api *APIRouter) services(resp http.ResponseWriter, r *http.Request) {
	var err error
	var routeList []byte
	var response common.APIResponse

	respType := r.URL.Query().Get("type")
	if respType == "full" || respType == "" {
		routes := api.serviceManager.GetAllServices()
		routeList, err = json.Marshal(routes)
	} else if respType == "minimal" {
		routes := api.serviceManager.GetServiceNames()
		routeList, err = json.Marshal(routes)
	}
	if err == nil {
		response = common.CreateAPIResponse(string(routeList), nil, 500)
	} else {
		response = common.CreateAPIResponse("failed", err, 500)
	}
	common.WriteAPIResponseStruct(resp, response)
}
func (api *APIRouter) process(resp http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	action := r.URL.Query().Get("action")
	switch action {
	case "start":
		if api.serviceManager.StartManagedService(name) {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 400))
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("", errors.New("not found or already running"), 400))
		}
	case "stop":
		api.serviceManager.StopManagedService(name)
	}
}
func (api *APIRouter) ws(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) wsauth(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) updateService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, api.serviceManager.UpdateService(r))
}
func (api *APIRouter) getService(resp http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if service, err := api.data.GetRoute(serviceName); err == nil {
		resBytes, _ := json.Marshal(service)
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(string(resBytes), nil, 404))
	} else {
		common.WriteFailureResponse(err, resp, "getService", 404)
	}
}
