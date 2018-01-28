package management

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.m/watchdog/common"
	"git.m/watchdog/data"
	"github.com/husobee/vestigo"
	"io/ioutil"
	"net/http"
)

//APIRouter ...
type APIRouter struct {
	user           *User
	router         *vestigo.Router
	serviceManager *ServiceManager
}

func NewAPIRouter(store *data.DataStore) *APIRouter {
	user := NewUserService(store)
	return &APIRouter{
		user:           user,
		serviceManager: NewServiceManager(store),
	}
}

//StartManagementAPIListener ...
func (api *APIRouter) StartManagementAPIListener() {
	api.router = vestigo.NewRouter()
	api.router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Cache-Control", "X-Requested-With", "Content-Type"},
		AllowOrigin:  []string{"https://watchdog.m.rdro.us", "http://192.168.1.12:4200"},
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

//InitAPIRouter ...
func (api *APIRouter) SetupRoutes() {
	api.router.Handle("/ws/log", common.RequestWrapper(common.Nothing, "GET", api.ws))
	api.router.Handle("/api/frost/wsauth", common.RequestWrapper(api.user.AuthTokenProvided, "POST", api.wsauth))

	api.router.Handle("/api/frost/user/login", common.RequestWrapper(common.Nothing, "POST", api.login))
	api.router.Handle("/api/frost/user/register", common.RequestWrapper(common.Nothing, "POST", api.register))

	api.router.Handle("/api/frost/service/new", common.RequestWrapper(api.user.AuthTokenProvided, "POST", api.newService))
	api.router.Handle("/api/frost/service/delete", common.RequestWrapper(api.user.AuthTokenProvided, "DELETE", api.deleteService))
	api.router.Handle("/api/frost/service/update", common.RequestWrapper(api.user.AuthTokenProvided, "PUT", api.updateService))

	api.router.Handle("/api/frost/services", common.RequestWrapper(api.user.AuthTokenProvided, "GET", api.services))

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
	var newService data.KnownRoute
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &newService)
	common.WriteAPIResponseStruct(resp, api.serviceManager.AddService(newService))
}
func (api *APIRouter) deleteService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) updateService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) services(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) ws(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) wsauth(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
