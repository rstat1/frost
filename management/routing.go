package management

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/husobee/vestigo"
	"github.com/rstat1/frost/auth"
	"github.com/rstat1/frost/common"
	"github.com/rstat1/frost/crypto"
	"github.com/rstat1/frost/data"
	"github.com/rstat1/frost/proxy"
	"github.com/rstat1/frost/services"
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
	vault          *crypto.VaultClient
}

//NewAPIRouter ...
func NewAPIRouter(store *data.DataStore, proxy *proxy.Proxy, serviceMan *ServiceManager, user *auth.User, vc *crypto.VaultClient, devMode bool) *APIRouter {
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
		vault:         vc,
		bingBGService: services.NewBingBGFetcher(store),
	}
}

//StartManagementAPIListener ...
func (api *APIRouter) StartManagementAPIListener() {
	api.serviceID, api.serviceKey = api.data.GetInstanceDetails()
	if api.dev {
		api.watchdog = "http://console" + common.BaseURL //"http://192.168.1.12:4200"
		api.baseAPIURL = "http://api" + common.BaseURL
		api.authServiceURL = "http://trinity" + common.BaseURL
	} else {
		api.watchdog = "https://console" + common.BaseURL
		api.baseAPIURL = "https://api" + common.BaseURL
		api.authServiceURL = "https://trinity" + common.BaseURL
	}
	api.router = vestigo.NewRouter()
	api.router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Cache-Control", "X-Requested-With", "Content-Type"},
		AllowOrigin: []string{"https://console" + common.BaseURL, "http://console" + common.BaseURL, "http://trinity" + common.BaseURL,
			"https://trinity" + common.BaseURL, "http://192.168.1.12:4200"},
	})
	vestigo.CustomNotFoundHandlerFunc(api.NotFound)
	api.SetupRoutes()

	if svcDetails, err := api.data.GetServiceByID(api.serviceID); err == nil {
		svcDetails.RedirectURL = api.watchdog + "/auth"
		common.LogError("", api.data.UpdateRoute(svcDetails, "watchdog"))
	}

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
	api.router.Handle("/api/frost/ws/log", common.RequestWrapper(common.Nothing, "GET", api.ws))
	api.router.Handle("/api/frost/wsauth", common.RequestWrapper(api.user.IsRoot, "POST", api.wsauth))

	api.router.Handle("/api/frost/auth/token", common.RequestWrapper(common.Nothing, "GET", api.getToken))

	api.router.Handle("/api/frost/service/get", common.RequestWrapper(api.user.IsRoot, "GET", api.getService))
	api.router.Handle("/api/frost/service/new", common.RequestWrapper(api.user.IsRoot, "POST", api.newService))
	api.router.Handle("/api/frost/service/edit", common.RequestWrapper(api.user.IsRoot, "POST", api.editService))
	api.router.Handle("/api/frost/service/delete", common.RequestWrapper(api.user.IsRoot, "DELETE", api.deleteService))
	api.router.Handle("/api/frost/service/update", common.RequestWrapper(api.user.IsRoot, "POST", api.updateService))
	api.router.Handle("/api/frost/service/restart/:name", common.RequestWrapper(api.user.IsRoot, "GET", api.restartService))

	api.router.Handle("/api/frost/aliases/new", common.RequestWrapper(api.user.IsRoot, "POST", api.newExtraRoute))
	api.router.Handle("/api/frost/aliases/all", common.RequestWrapper(api.user.IsRoot, "GET", api.getExtraRoutes))
	api.router.Handle("/api/frost/aliases/delete", common.RequestWrapper(api.user.IsRoot, "POST", api.deleteExtraRoute))

	// api.router.Handle("/api/frost/reverseproxy/addroute", common.RequestWrapper(api.user.IsRoot, "POST", api.newproxyroute))
	// api.router.Handle("/api/frost/reverseproxy/delteroute", common.RequestWrapper(api.user.IsRoot, "DELETE", api.deleteproxyroute))

	api.router.Handle("/api/frost/services", common.RequestWrapper(api.user.IsRoot, "GET", api.services))

	api.router.Handle("/api/frost/process", common.RequestWrapper(api.user.IsRoot, "GET", api.process))
	api.router.Handle("/api/frost/bg", common.RequestWrapper(common.Nothing, "GET", api.bingBGService.GetBGImage))

	api.router.Handle("/api/frost/icons", common.RequestWrapper(common.Nothing, "GET", api.icons))
	api.router.Handle("/api/frost/icon/:service", common.RequestWrapper(common.Nothing, "GET", api.geticon))
	api.router.Handle("/api/frost/icon/new/:service", common.RequestWrapper(api.user.IsRoot, "POST", api.newicon))

	api.router.Handle("/api/frost/init", common.RequestWrapper(common.Nothing, "GET", api.initFrost))
	api.router.Handle("/api/frost/status", common.RequestWrapper(common.Nothing, "GET", api.firstRunStatus))
	api.router.Handle("/api/frost/serviceid", common.RequestWrapper(common.Nothing, "GET", api.getServiceID))

}
func (api *APIRouter) newExtraRoute(resp http.ResponseWriter, r *http.Request) {
	var propChange data.RouteAlias
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &propChange)
		common.WriteAPIResponseStruct(resp, api.servMan.AddNewExtraRoute(propChange))
	} else {
		common.WriteFailureResponse(err, resp, "newExtraRoute", 500)
	}
}
func (api *APIRouter) getExtraRoutes(resp http.ResponseWriter, r *http.Request) {
	apiName := r.URL.Query().Get("api")
	if apiName != "" {
		common.WriteAPIResponseStruct(resp, api.servMan.GetExtraRoutes(apiName))
	} else {
		common.WriteFailureResponse(errors.New("specify service API name"), resp, "getExtraRoutes", 400)
	}
}
func (api *APIRouter) deleteExtraRoute(resp http.ResponseWriter, r *http.Request) {
	var deleteRequest data.AliasDeleteRequest
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &deleteRequest)
		api.proxy.DeleteExtraRoute(deleteRequest.BaseURL, deleteRequest.Route)
		apiResp := common.CreateAPIResponse("success", api.data.DeleteExtraRoute(deleteRequest.Route), 500)
		common.WriteAPIResponseStruct(resp, apiResp)
	}
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
	var oldServiceName, oldBinName string
	var e error
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		json.Unmarshal(body, &propChange)
		if propChange.PropertyName != "newroute" {
			if service, err := api.data.GetRoute(propChange.ServiceName); err == nil {
				switch propChange.PropertyName {
				case "name":
					oldServiceName = service.AppName
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
				case "filename":
					oldBinName = service.BinName
					service.BinName = propChange.NewValue
					break
				case "vault":
					if propChange.NewValue == "Enabled" {
						service.VaultIntegrated = true
					} else {
						service.VaultIntegrated = false
					}
					break
				}
				if propChange.ServiceName == "watchdog" {
					e = api.data.UpdateSysConfig(propChange)
				}
				e = api.data.UpdateRoute(service, propChange.ServiceName)
				if propChange.PropertyName == "filename" {
					e = api.servMan.RenameServiceBin(oldBinName, service.BinName, service.AppName)
				}
				if propChange.PropertyName == "name" {
					e = api.servMan.RenameServiceDirectory(oldServiceName, propChange.NewValue)
					api.proxy.RenameRoute(oldServiceName, propChange.NewValue)
				}
			} else {
				e = err
			}
		} else {
			newRoute := data.ExtraRoute{
				APIName: propChange.ServiceName,
				FullURL: propChange.NewValue,
			}
			e = api.data.AddExtraRoute(newRoute)
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
func (api *APIRouter) restartService(resp http.ResponseWriter, r *http.Request) {
	serviceName := vestigo.Param(r, "name")
	if serviceName != "" {
		api.servMan.StopManagedService(serviceName)
		if api.servMan.StartManagedService(serviceName) == false {
			common.WriteFailureResponse(errors.New("failed to start service: "+serviceName), resp, "restartService", 500)
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 200))
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", errors.New("no service specified"), 400))
	}
}
func (api *APIRouter) firstRunStatus(resp http.ResponseWriter, r *http.Request) {
	if !api.data.GetFirstRunState() {
		if api.vault.TokenSet {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("initialized", nil, 200))
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("initialized-need-vt", nil, 200))
		}
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
		common.LogInfo("password", password, "First run init complete. Restart Frost, and login with the username 'root' and the provided password, to continue.")
		os.Exit(0)
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
func (api *APIRouter) newicon(resp http.ResponseWriter, r *http.Request) {
	var service = vestigo.Param(r, "service")
	if err := r.ParseMultipartForm(1 * 1024 * 1024); err == nil {
		if icon, _, e := r.FormFile("icon"); e == nil {
			if e = api.servMan.handleIconFileUpload(icon, service); e != nil {
				common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", err, 500))
			}
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 200))
		}
	}
}
func (api *APIRouter) geticon(resp http.ResponseWriter, r *http.Request) {
	var iconPath string
	var service = vestigo.Param(r, "service")
	if _, e := os.Stat("console/serviceicons/" + service + ".png"); e != nil {
		iconPath = "console/web/assets/services.png"
		common.Logger.Errorln(e)
	} else {
		iconPath = "console/serviceicons/" + service + ".png"
	}
	if icon, e := os.Open(iconPath); e == nil {
		if image, err := ioutil.ReadAll(icon); err == nil {
			resp.Write(image)
		}
	} else {
		common.Logger.Errorln(e)
		resp.WriteHeader(404)
	}
}
func (api *APIRouter) icons(resp http.ResponseWriter, r *http.Request) {
	var icons []string
	files, _ := ioutil.ReadDir("console/serviceicons/")
	for _, file := range files {
		icons = append(icons, file.Name())
	}
	f, _ := json.Marshal(icons)
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(string(f), nil, 200))
}
func (api *APIRouter) newproxyroute(resp http.ResponseWriter, r *http.Request) {
	var proxiedRoute data.ProxyRoute
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, &proxiedRoute); err != nil {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", err, 500))
		} else {
			apiResp := common.CreateAPIResponse("success", api.data.AddProxyRoute(proxiedRoute.Hostname, proxiedRoute.IPAddress), 500)
			common.WriteAPIResponseStruct(resp, apiResp)
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", err, 500))
	}
}
func (api *APIRouter) deleteproxyroute(resp http.ResponseWriter, r *http.Request) {

}
