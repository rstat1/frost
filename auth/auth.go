package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"github.com/husobee/vestigo"
)

const (
	requestIDTTL = 120

	prodLoginURL = "https://trinity.m.rdro.us"
	devLoginURL  = "http://trinity.dev-m.rdro.us" //"http://192.168.1.12:4200"
)

//AuthService ...
type AuthService struct {
	user       *User
	db         *data.DataStore
	route      *vestigo.Router
	cache      *data.CacheService
	loginURL   string
	inDevMode  bool
	serviceID  string
	serviceKey string
	httpClient *http.Client
}

//NewAuthService ...
func NewAuthService(db *data.DataStore, user *User, devmode bool) *AuthService {
	var http = &http.Client{
		Timeout: time.Second * 2,
	}
	id, key := db.GetInstanceDetails()

	return &AuthService{
		db:         db,
		user:       user,
		cache:      db.Cache,
		inDevMode:  devmode,
		httpClient: http,
		serviceID:  id,
		serviceKey: key,
	}
}

//InitAuthService ...
func (auth *AuthService) InitAuthService() {
	if auth.inDevMode {
		auth.loginURL = devLoginURL
	} else {
		auth.loginURL = prodLoginURL
	}
	auth.route = vestigo.NewRouter()
	auth.route.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Cache-Control", "X-Requested-With", "Content-Type"},
		AllowOrigin:  []string{"*"},
	})
	auth.initAPIRoutes()
	go func() {
		if err := http.ListenAndServe("localhost:1003", auth.route); err != nil {
			common.CreateFailureResponse(err, "StartManagementAPIListener", 500)
		}
	}()
}
func (auth *AuthService) initAPIRoutes() {
	if auth.inDevMode {
		auth.loginURL = devLoginURL
	} else {
		auth.loginURL = prodLoginURL
	}
	auth.route.Handle("/api/trinity/", common.RequestWrapper(auth.user.AuthTokenProvided, "GET", auth.echo))

	auth.route.Handle("/api/trinity/token", common.RequestWrapper(auth.CodeAndKeyProvided, "GET", auth.token))
	auth.route.Handle("/api/trinity/validate", common.RequestWrapper(auth.CredsAndIDProvided, "POST", auth.validate))
	auth.route.Handle("/api/trinity/authorize", common.RequestWrapper(auth.HasServiceID, "GET", auth.authorize))

	auth.route.Handle("/api/trinity/ws/ticket", common.RequestWrapper(auth.user.AuthTokenProvided, "POST", auth.wsticket))
	auth.route.Handle("/api/trinity/ws/validate", common.RequestWrapper(common.Nothing, "GET", auth.checkwsticket))

	auth.route.Handle("/api/trinity/permissions", common.RequestWrapper(auth.user.IsRoot, "GET", auth.permissions))
	auth.route.Handle("/api/trinity/permissions/change", common.RequestWrapper(auth.user.IsRoot, "POST", auth.changepermission))

	auth.route.Handle("/api/trinity/user", common.RequestWrapper(auth.user.AuthTokenProvided, "GET", auth.userinfo))
	auth.route.Handle("/api/trinity/user/new", common.RequestWrapper(auth.user.AuthTokenProvided, "POST", auth.newuser))
	auth.route.Handle("/api/trinity/user/list", common.RequestWrapper(auth.user.AuthTokenProvided, "GET", auth.getusers))

	auth.route.Handle("/api/trinity/user/edit", common.RequestWrapper(auth.user.IsRoot, "POST", auth.edituser))
	auth.route.Handle("/api/trinity/user/delete", common.RequestWrapper(auth.user.IsRoot, "DELETE", auth.deleteuser))

	auth.route.Handle("/api/trinity/service/fromrid", common.RequestWrapper(auth.HasRequestID, "GET", auth.fromrequest))
}

//NotFound ...
func (auth *AuthService) NotFound(resp http.ResponseWriter, r *http.Request) {
	common.WriteFailureResponse(errors.New("route "+r.URL.String()+" not found."), resp, "NotFound", 404)
}

//HasServiceID ...
func (auth *AuthService) HasServiceID(r *http.Request) common.APIResponse {
	var resp common.APIResponse
	serviceID := auth.hasRequiredParam("sid", r)

	if serviceID.Status == "success" {
		if _, err := auth.db.GetServiceByID(serviceID.Response); err == nil {
			resp = serviceID
		} else {
			if serviceID.Response == auth.serviceID {
				resp = serviceID
			} else {
				resp = common.CreateFailureResponse(errors.New("unknown service id"), "HasServiceID", 400)
			}
		}
	} else {
		resp = common.CreateFailureResponse(errors.New("missing required parameter"), "HasServiceID", 400)
	}
	return resp
}

//HasServiceCreds ...
func (auth *AuthService) HasServiceCreds(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	hasSid := auth.hasRequiredParam("sid", r)
	hasSKey := auth.hasRequiredParam("skey", r)
	if hasSid.Status == "success" && hasSKey.Status == "success" {
		resp = hasSid
	} else {
		resp = common.CreateFailureResponse(errors.New("missing required parameter"), "HasServiceCreds", 400)
	}
	return resp
}

//CodeAndKeyProvided ...
func (auth *AuthService) CodeAndKeyProvided(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	hasSid := auth.hasRequiredParam("sid", r)
	hasSKey := auth.hasRequiredParam("skey", r)
	hasCode := auth.hasRequiredParam("code", r)
	if hasSid.Status == "success" && hasSKey.Status == "success" && hasCode.Status == "success" {
		resp = hasSid
	} else {
		resp = common.CreateFailureResponse(errors.New("missing required parameter(s)"), "HasServiceCreds", 400)
	}
	return resp
}

//CredsAndIDProvided ...
func (auth *AuthService) CredsAndIDProvided(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	requestID := auth.hasRequiredParam("r", r)
	serviceID := auth.cache.GetString("authrequest", requestID.Response)

	if requestID.Status == "success" && serviceID != "" && r.ContentLength > 0 {
		resp = common.CreateAPIResponse("success", nil, 400)
	} else {
		resp = common.CreateFailureResponse(errors.New("missing required parameter(s)"), "CredsAndIDProvided", 400)
	}
	return resp
}

//HasRequestID ...
func (auth *AuthService) HasRequestID(r *http.Request) common.APIResponse {
	var resp common.APIResponse

	requestID := auth.hasRequiredParam("r", r)
	if requestID.Status == "success" {
		sid := auth.cache.GetString("authrequest", requestID.Response)
		if sid == "" {
			resp = common.CreateAPIResponse("failed", errors.New("invalid request"), 500)
		} else {
			resp = requestID
		}
	} else {
		resp = requestID
	}
	return resp
}
func (auth *AuthService) hasRequiredParam(param string, r *http.Request) common.APIResponse {
	value := r.URL.Query().Get(param)
	if value == "" {
		return common.CreateFailureResponse(errors.New("missing required param"), "hasRequiredParam", 500)
	}
	return common.CreateAPIResponse(value, nil, 200)
}
func (auth *AuthService) token(resp http.ResponseWriter, r *http.Request) {
	hasSid := auth.hasRequiredParam("sid", r)
	hasSKey := auth.hasRequiredParam("skey", r)
	hasCode := auth.hasRequiredParam("code", r)
	if username := auth.cache.GetString("auth-"+hasSid.Response, hasCode.Response); username != "" {
		if service, err := auth.db.GetServiceByID(hasSid.Response); err == nil {
			if service.ServiceKey == hasSKey.Response {
				common.WriteAPIResponseStruct(resp, auth.user.GenerateAuthToken(username, service.AppName))
			} else {
				common.WriteFailureResponse(errors.New("unknown service key"), resp, "token", 400)
			}
		} else {
			common.WriteFailureResponse(err, resp, "token", 500)
		}
	} else {
		common.WriteFailureResponse(errors.New("auth code expired or invalid"), resp, "token", 500)
	}
}

func (auth *AuthService) userinfo(resp http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("name")
	var user data.User

	if username == "" {
		user = auth.user.GetUserFromToken(r)
	} else {
		user = auth.db.GetUser(username)
	}

	if u, err := json.Marshal(user); err == nil {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(string(u), nil, 200))
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", err, 200))
	}
}
func (auth *AuthService) validate(resp http.ResponseWriter, r *http.Request) {
	var authRequest data.AuthRequest
	requestID := auth.hasRequiredParam("r", r)
	sid := auth.cache.GetString("authrequest", requestID.Response)
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &authRequest)
	if service, err := auth.db.GetServiceByID(sid); err == nil {
		if response := auth.user.ValidateLoginRequest(authRequest); response.Status == "success" {
			if auth.db.DoesUserHavePermission(authRequest.Username, service.AppName, "hasAccess") {
				authCode := common.RandomID(48)
				redirect := service.RedirectURL + "?type=authcode&code=" + authCode
				auth.cache.PutStringWithExpiration("auth-"+service.ServiceID, authCode, authRequest.Username, requestIDTTL*5)
				common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(redirect, nil, 500))
			} else {
				common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", errors.New("not authorized"), 401))
			}
		} else {
			common.WriteAPIResponseStruct(resp, response)
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", err, 400))
		common.CreateFailureResponse(err, "validate", 500)
	}
}
func (auth *AuthService) authorize(resp http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("sid")
	requestID := common.RandomID(48)
	auth.cache.PutStringWithExpiration("authrequest", requestID, service, requestIDTTL)
	resp.Header().Set("Location", auth.loginURL+"/login?r="+requestID)
	resp.WriteHeader(303)
}
func (auth *AuthService) newuser(resp http.ResponseWriter, r *http.Request) {
	var request data.UserDetails
	if strings.Contains(r.Host, "localhost") {
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &request); err == nil {
			auth.user.ParsePermissionList(request.Permissions)
			common.WriteAPIResponseStruct(resp, auth.user.NewUser(data.AuthRequest{
				Username: request.Username,
				Password: request.Password,
			}, request.Permissions))
		} else {
			common.WriteFailureResponse(fmt.Errorf("failed deserializing request body %s", err), resp, "register", 500)
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("failed", errors.New("invalid request"), 400))
	}
}
func (auth *AuthService) fromrequest(resp http.ResponseWriter, r *http.Request) {
	var response common.APIResponse
	value := r.URL.Query().Get("r")
	sid := auth.cache.GetString("authrequest", value)

	if s, err := auth.db.GetServiceByID(sid); err == nil {
		response = common.CreateAPIResponse(s.AppName, nil, 200)
	} else {
		if sid == auth.serviceID {
			response = common.CreateAPIResponse("watchdog", nil, 200)
		} else {
			response = common.CreateFailureResponse(err, "fromrequest", 500)
		}
	}
	common.WriteAPIResponseStruct(resp, response)
}
func (auth *AuthService) echo(resp http.ResponseWriter, r *http.Request) {
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", nil, 500))
}
func (auth *AuthService) getusers(resp http.ResponseWriter, r *http.Request) {
	if users, err := auth.db.GetAllUserNames(); err == nil {
		if list, e := json.Marshal(users); e == nil {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(string(list), nil, 200))
		} else {
			common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("", e, 500))
		}
	} else {
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("", err, 500))
	}
}
func (auth *AuthService) deleteuser(resp http.ResponseWriter, r *http.Request) {
	var name = r.URL.Query().Get("name")
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", auth.db.DeleteUser(name), 400))
}
func (auth *AuthService) edituser(resp http.ResponseWriter, r *http.Request) {
	var request data.PasswordChange
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err == nil {
		passHasher := sha512.New512_256()
		hash := passHasher.Sum([]byte(request.Password))
		user := auth.db.GetUser(request.Username)
		user.PassHash = hex.EncodeToString(hash)
		common.WriteAPIResponseStruct(resp, common.CreateAPIResponse("success", auth.db.UpdateUser(user), 400))
	} else {
		common.WriteFailureResponse(err, resp, "edituser", 500)
	}
}
func (auth *AuthService) permissions(writer http.ResponseWriter, r *http.Request) {
	var resp common.APIResponse
	var name = r.URL.Query().Get("user")
	if name != "" {
		permMap, err := auth.db.GetUserPermissionMap(name)
		if err == nil {
			if list, e := json.Marshal(permMap); e == nil {
				resp = common.CreateAPIResponse(string(list), e, 500)
			} else {
				resp = common.CreateAPIResponse("failed", e, 500)
			}
		} else {
			resp = common.CreateAPIResponse("failed", err, 500)
		}
		common.WriteAPIResponseStruct(writer, resp)
	} else {
		common.WriteFailureResponse(errors.New("no username specified"), writer, "permission", 400)
	}
}
func (auth *AuthService) changepermission(resp http.ResponseWriter, r *http.Request) {
	var request data.PermissionChange

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err == nil {
		common.Logger.Debugln(request)
		common.Logger.Errorln(auth.db.UpdateUserPermissions(request))
	}
}
func (auth *AuthService) checkwsticket(resp http.ResponseWriter, r *http.Request) {
	var apiResp common.APIResponse
	if t := r.URL.Query().Get("t"); t != "" {
		userID := auth.cache.GetString(t, "userid")
		if userID != "" {
			auth.cache.DeleteString(t, "userid")
			apiResp = common.CreateAPIResponse(userID, nil, 200)
		} else {
			apiResp = common.CreateAPIResponse("", errors.New("no userid"), 200)
		}
	} else {
		apiResp = common.CreateFailureResponse(errors.New("no ticket"), "checkwsticket", 400)
	}
	common.WriteAPIResponseStruct(resp, apiResp)
}
func (auth *AuthService) wsticket(resp http.ResponseWriter, r *http.Request) {
	ticket := common.RandomID(48)
	user := auth.user.GetUserFromToken(r)
	auth.cache.PutString(ticket, "userid", user.Id)
	common.WriteAPIResponseStruct(resp, common.CreateAPIResponse(ticket, nil, 400))
}
