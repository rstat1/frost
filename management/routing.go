package management

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"git.m/svcman/auth"
	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/proxy"
	"git.m/svcman/services"
	"github.com/husobee/vestigo"
)

const (
	prodBaseURL = ".m.rdro.us"

	devBaseURL = ".dev-m.rdro.us"
)

//InstanceInfo ...
type InstanceInfo struct {
	ServiceID, ServiceKey, Password string
}

//APIRouter ...
type APIRouter struct {
	dev            bool
	user           *auth.User
	data           *data.DataStore
	proxy          *proxy.Proxy
	router         *vestigo.Router
	servMan        *ServiceManager
	bingBGService  *services.BingBGFetcher
	authServiceURL string
	baseAPIURL     string
	serviceKey     string
	serviceID      string
	watchdog       string
	httpClient     *http.Client
}

//NewAPIRouter ...
func NewAPIRouter(store *data.DataStore, proxy *proxy.Proxy, serviceMan *ServiceManager, user *auth.User, devMode bool) *APIRouter {
	var http = &http.Client{
		Timeout: time.Second * 2,
	}

	return &APIRouter{
		dev:           devMode,
		user:          user,
		data:          store,
		proxy:         proxy,
		servMan:       serviceMan,
		httpClient:    http,
		bingBGService: services.NewBingBGFetcher(store),
	}
}

//StartManagementAPIListener ...
func (api *APIRouter) StartManagementAPIListener() {
	api.serviceID, api.serviceKey = api.data.GetInstanceDetails()
	if api.dev {
		api.watchdog = "http://watchdog" + devBaseURL //"http://192.168.1.12:4200"
		api.baseAPIURL = "http://api" + devBaseURL
		api.authServiceURL = "http://trinity" + devBaseURL
	} else {
		api.watchdog = "https://watchdog" + prodBaseURL
		api.baseAPIURL = "https://api" + prodBaseURL
		api.authServiceURL = "https://trinity" + prodBaseURL
	}
	api.router = vestigo.NewRouter()
	api.router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Cache-Control", "X-Requested-With", "Content-Type"},
		AllowOrigin: []string{"https://watchdog.m.rdro.us", "http://trinity.dev-m.rdro.us", "https://trinity.m.rdro.us",
			"http://192.168.1.12:4200", "http://watchdog.dev-m.rdro.us"},
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

	api.router.Handle("/api/frost/auth/token", common.RequestWrapper(common.Nothing, "GET", api.getToken))

	api.router.Handle("/api/frost/service/get", common.RequestWrapper(api.user.IsRoot, "GET", api.getService))
	api.router.Handle("/api/frost/service/new", common.RequestWrapper(api.user.IsRoot, "POST", api.newService))
	api.router.Handle("/api/frost/service/edit", common.RequestWrapper(api.user.IsRoot, "POST", api.editService))
	api.router.Handle("/api/frost/service/delete", common.RequestWrapper(api.user.IsRoot, "DELETE", api.deleteService))
	api.router.Handle("/api/frost/service/update", common.RequestWrapper(api.user.IsRoot, "POST", api.updateService))

	api.router.Handle("/api/frost/services", common.RequestWrapper(api.user.IsRoot, "GET", api.services))

	api.router.Handle("/api/frost/process", common.RequestWrapper(api.user.IsRoot, "GET", api.process))
	api.router.Handle("/api/frost/bg", common.RequestWrapper(common.Nothing, "GET", api.bingBGService.GetBGImage))

	api.router.Handle("/api/frost/status", common.RequestWrapper(common.Nothing, "GET", api.firstRunStatus))
	api.router.Handle("/api/frost/init", common.RequestWrapper(common.Nothing, "GET", api.initFrost))
	api.router.Handle("/api/frost/serviceid", common.RequestWrapper(common.Nothing, "GET", api.getServiceID))
}
func (api *APIRouter) getToken(resp http.ResponseWriter, r *http.Request) {
	var serviceResp common.APIResponse

	code := r.URL.Query().Get("code")
	req, _ := http.NewRequest("GET", api.baseAPIURL+"/trinity/token", nil)
	q := req.URL.Query()
	q.Set("sid", api.serviceID)
	q.Set("skey", api.serviceKey)
	q.Set("code", code)
	req.URL.RawQuery = q.Encode()
	if httpResp, err := api.httpClient.Do(req); err == nil {
		if body, err := ioutil.ReadAll(httpResp.Body); err != nil {
			common.WriteFailureResponse(err, resp, "getToken", 500)
		} else {
			if e := json.Unmarshal(body, &serviceResp); e != nil {
				common.WriteFailureResponse(e, resp, "getToken", 500)
			} else {
				if serviceResp.Status == "failed" {
					serviceResp.HttpStatusCode = 500
				} else {
					serviceResp.HttpStatusCode = 200
				}
				common.WriteAPIResponseStruct(resp, serviceResp)
			}
		}
	} else {
		common.WriteFailureResponse(err, resp, "getToken", 500)
	}
}
func (api *APIRouter) newService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, api.getServiceListOnSuccess(api.servMan.NewService(r)))
}
func (api *APIRouter) deleteService(resp http.ResponseWriter, r *http.Request) {
	var name = r.URL.Query().Get("name")
	api.servMan.StopManagedService(name)
	common.WriteAPIResponseStruct(resp, api.getServiceListOnSuccess(api.servMan.DeleteService(name)))
}
func (api *APIRouter) editService(resp http.ResponseWriter, r *http.Request) {
	var propChange data.ServiceEdit
	var e error
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &propChange)
		if service, err := api.data.GetRoute(propChange.ServiceName); err == nil {
			switch propChange.PropertyName {
			case "name":
				service.AppName = propChange.NewValue
				break
			case "apiName":
				service.APIName = propChange.NewValue
				break
			case "skey":
				propChange.NewValue = common.RandomID(48)
				if propChange.ServiceName == "watchdog" {
					api.serviceKey = propChange.NewValue
				}
				service.ServiceKey = propChange.NewValue
				break
			case "redirect":
				service.RedirectURL = propChange.NewValue
				break
			case "localaddr":
				service.ServiceAddress = propChange.NewValue
				break
			case "managed":
				if propChange.NewValue == "Enabled" {
					service.IsManagedService = true
				} else {
					service.IsManagedService = false
				}
				break
			}
			if propChange.ServiceName == "watchdog" {
				e = api.data.UpdateSysConfig(propChange)
			}
			e = api.data.UpdateRoute(service, propChange.ServiceName)
		} else {
			e = err
		}
	} else {
		e = err
	}
	if e != nil {
		common.WriteFailureResponse(e, resp, "editService", 400)
	} else {
		if propChange.PropertyName == "skey" {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(propChange.NewValue, nil, 200))
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 200))
		}
	}
}
func (api *APIRouter) services(resp http.ResponseWriter, r *http.Request) {
	var err error
	var routeList []byte
	var response common.APIResponse

	respType := r.URL.Query().Get("type")
	if respType == "full" || respType == "" {
		routes := api.servMan.GetAllServices()
		routeList, err = json.Marshal(routes)
	} else if respType == "minimal" {
		routes := api.servMan.GetServiceNames()
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
		if api.servMan.StartManagedService(name) {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 400))
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("", errors.New("not found or already running"), 400))
		}
	case "stop":
		api.servMan.StopManagedService(name)
	}
}
func (api *APIRouter) ws(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) wsauth(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("not implemented", nil, 501))
}
func (api *APIRouter) updateService(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, api.servMan.UpdateService(r))
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
func (api *APIRouter) firstRunStatus(resp http.ResponseWriter, r *http.Request) {
	if !api.data.GetFirstRunState() {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("initialized", nil, 200))
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(api.watchdog+"/first-run", nil, 200))
	}
}
func (api *APIRouter) initFrost(resp http.ResponseWriter, r *http.Request) {
	var perms []data.ServiceAuth
	if api.data.GetFirstRunState() == true {
		sid, skey := api.data.GetInstanceDetails()
		password := common.RandomID(48)
		newUser := data.AuthRequest{
			Username: "root",
			Password: password,
		}
		service := data.ServiceDetails{
			AppName:          "watchdog",
			BinName:          "watchdog",
			APIName:          "frost",
			IsManagedService: false,
			ServiceID:        sid,
			ServiceKey:       skey,
			RedirectURL:      api.watchdog + "/auth",
			ServiceAddress:   "localhost:1000",
		}
		p := data.ServiceAuth{
			Service: "watchdog",
			Permissions: []data.PermissionValue{
				{Name: "hasAccess", Value: true},
				{Name: "hasRoot", Value: true},
			},
		}
		perms := append(perms, p)
		if err := api.data.AddNewRoute(service); err != nil {
			common.Logger.WithField("func", "initFrost").Errorln(err)
		}
		api.user.NewUser(newUser, perms)
		api.data.SetFirstRunState()
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(password, nil, 200))
	} else {
		common.WriteFailureResponse(errors.New("already initialized"), resp, "initFrost", 400)
	}
}
func (api *APIRouter) getServiceID(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(api.serviceID, nil, 400))
}
func (api *APIRouter) getServiceListOnSuccess(resp common.APIResponse) common.APIResponse {
	var apiResp common.APIResponse
	if resp.Status == "success" {
		routes := api.servMan.GetServiceNames()
		routeList, _ := json.Marshal(routes)
		apiResp = common.CreateAPIResponse(string(routeList), nil, 500)
	} else {
		apiResp = resp
	}
	return apiResp
}
