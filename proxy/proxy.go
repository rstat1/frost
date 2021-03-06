package proxy

import (
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"time"

	"net/url"
	"strings"

	"github.com/rstat1/frost/common"
	"github.com/rstat1/frost/data"
	"github.com/rstat1/frost/services"

	"github.com/pkg/errors"
	"github.com/vulcand/oxy/forward"
)

//Proxy ...
type Proxy struct {
	data                 *data.DataStore
	fwd                  *forward.Forwarder
	internal             *InternalProxy
	knownRoutes          map[string]bool
	apiRoutes            map[string]string
	aliasHosts           map[string][]string
	knownAliasHosts      map[string]bool
	hostsToAPIServer     map[string]string
	reverseProxyRoutes   map[string]string
	webuiserver          *services.WebServer
	isInLocalMode        bool
	icAPIURL             string
	apiBaseURL           string
	baseURL              string
	baseAuthURL          string
	baseWDURL            string
	apiBaseURLWithScheme string
	listenerPort         string
	apiNameToOrigin      map[string]string
	botFun               [20]string

	internalBaseURL           string
	internalAPIBase           string
	internalAPIBaseWithScheme string
	internalRoutes            map[string]string
}

var (
	lastHost, lastHostPort string
)

const (
	watchdogAPIName    = "frost"
	authServiceAPIName = "trinity"

	devPort        = ":80"
	productionPort = ":443"
)

//HTTPErrorHandler This only exists because vulcand/oxy doesn't implment handling request errors in a way that any reasonable person would consider sane.
type HTTPErrorHandler struct{}

func (sieh *HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, err error) {
	common.LogError(req.URL.String(), err)
	common.WriteFailureResponse(err, w, "whoKnows", 500)
}

//NewProxy ...
func NewProxy(dataStoreRef *data.DataStore, devMode *bool) *Proxy {
	var fwd *forward.Forwarder
	sieh := &HTTPErrorHandler{}
	fwd, _ = forward.New(forward.ErrorHandler(sieh))
	p := &Proxy{
		fwd:                fwd,
		data:               dataStoreRef,
		internal:           NewInternalProxy(devMode),
		apiRoutes:          make(map[string]string),
		aliasHosts:         make(map[string][]string),
		knownRoutes:        make(map[string]bool),
		isInLocalMode:      *devMode,
		knownAliasHosts:    make(map[string]bool),
		webuiserver:        services.NewWebServer(),
		hostsToAPIServer:   make(map[string]string),
		apiNameToOrigin:    make(map[string]string),
		internalRoutes:     make(map[string]string),
		reverseProxyRoutes: make(map[string]string),
	}

	p.botFun = [20]string{
		"https://www.youtube.com/watch?v=wbby9coDRCk",
		"https://www.youtube.com/watch?v=nb2evY0kmpQ",
		"https://www.youtube.com/watch?v=eh7lp9umG2I",
		"https://www.youtube.com/watch?v=z9Uz1icjwrM",
		"https://www.youtube.com/watch?v=Sagg08DrO5U",
		"https://www.youtube.com/watch?v=5XmjJvJTyx0",
		"https://www.youtube.com/watch?v=IkdmOVejUlI",
		"https://www.youtube.com/watch?v=jScuYd3_xdQ",
		"https://www.youtube.com/watch?v=S5PvBzDlZGs",
		"https://www.youtube.com/watch?v=9UZbGgXvCCA",
		"https://www.youtube.com/watch?v=O-dNDXUt1fg",
		"https://www.youtube.com/watch?v=MJ5JEhDy8nE",
		"https://www.youtube.com/watch?v=VnnWp_akOrE",
		"https://www.youtube.com/watch?v=jwGfwbsF4c4",
		"https://www.youtube.com/watch?v=8ZcmTl_1ER8",
		"https://www.youtube.com/watch?v=gLmcGkvJ-e0",
		"https://www.youtube.com/watch?v=hGlyFc79BUE",
		"https://www.youtube.com/watch?v=xA8-6X8aR3o",
		"https://www.youtube.com/watch?v=7R1nRxcICeE",
		"https://www.youtube.com/watch?v=sCNrK-n68CM",
	}

	return p
}

//StartProxyListener ...
func (p *Proxy) StartProxyListener() {
	common.Logger.WithField("localmode", p.isInLocalMode).Infoln("starting proxy listener.")

	p.baseURL = common.BaseURL
	p.internalBaseURL = "frost.m"
	p.apiBaseURL = "api" + p.baseURL
	p.baseAuthURL = "trinity" + p.baseURL
	p.baseWDURL = "console" + p.baseURL

	if p.isInLocalMode == false {
		p.listenerPort = productionPort
		p.apiBaseURLWithScheme = "https://" + "." + p.baseURL
		p.internalAPIBase = "api.frost.m"
		p.internalAPIBaseWithScheme = "https//" + ".frost.m"
		common.Logger.Infoln("running in production mode...")
		p.setRoutes()
		go p.internal.StartProxyListener()
		p.startTLSServer()
	} else {
		p.listenerPort = devPort
		p.apiBaseURLWithScheme = "http://" + "." + p.baseURL
		p.internalAPIBase = "api.frost-int.m"
		p.internalAPIBaseWithScheme = "http://" + ".frost-int.m"
		common.Logger.Infoln("running in dev mode...")
		p.setRoutes()
		go p.internal.StartProxyListener()
		p.startNotTLSServer()
	}
}

//AddRoute ...
func (p *Proxy) AddRoute(newRoute data.ServiceDetails) {
	if newRoute.Internal {
		p.internalRoutes[newRoute.APIName] = newRoute.ServiceAddress
		p.knownRoutes[newRoute.AppName+p.internalBaseURL] = true
	} else {
		p.apiRoutes[newRoute.APIName] = newRoute.ServiceAddress
		p.knownRoutes[newRoute.AppName+p.baseURL] = true
	}
}

//AddExtraRoute ...
func (p *Proxy) AddExtraRoute(newExtraRoute data.ExtraRoute) {
	if strings.Contains(newExtraRoute.APIRoute, "*") {
		newExtraRoute.APIRoute = newExtraRoute.APIRoute[0 : len(newExtraRoute.APIRoute)-1]
	}
	p.aliasHosts[newExtraRoute.FullURL] = append(p.aliasHosts[newExtraRoute.FullURL], newExtraRoute.APIRoute)
	p.knownAliasHosts[newExtraRoute.FullURL] = true
	p.hostsToAPIServer[newExtraRoute.FullURL] = newExtraRoute.APIName
}

//DeleteRoute ...
func (p *Proxy) DeleteRoute(apiName, appName string) {
	delete(p.apiRoutes, apiName)
	delete(p.knownRoutes, appName+p.baseURL)
}

//DeleteExtraRoute ...
func (p *Proxy) DeleteExtraRoute(fullURL, apiRoute string) {
	delete(p.aliasHosts, fullURL)
	delete(p.knownAliasHosts, fullURL)
	delete(p.hostsToAPIServer, fullURL)
}

//RenameRoute ...
func (p *Proxy) RenameRoute(oldname, newname string) {
	p.knownRoutes = make(map[string]bool)
	p.apiRoutes = make(map[string]string)
	p.setRoutes()
}

//ServeHTTP ...
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// if val, ok := p.reverseProxyRoutes[req.Host]; ok {
	// 	p.proxyRequest(w, req, val, req.URL.Path, "", true)
	// } else {

	if strings.HasSuffix(req.URL.String(), ".php") {
		http.Redirect(w, req, p.botFun[rand.Intn(len(p.botFun))], 301)
		return
	}

	if req.URL.String() == "/favicon.ico" {
		w.WriteHeader(200)
	} else if req.Host == p.apiBaseURL {
		p.serveAPIRequest(w, req)
	} else if strings.HasSuffix(req.Host, p.baseURL) && p.knownAliasHosts[req.Host] == false {
		if p.knownRoutes[req.Host] {
			var name = strings.Replace(req.Host, p.baseURL, "", -1)
			p.webuiserver.ServeWebRequest(w, req, name)
		} else {
			common.WriteFailureResponse(errors.New("unknown route "+req.Host+req.URL.String()), w, "ServeHTTP", 404)
		}
	} else if p.knownAliasHosts[req.Host] == true {
		p.serveExtraRouteRequest(w, req)
	} else {
		common.WriteFailureResponse(errors.New("unknown route "+req.Host+req.URL.String()), w, "ServeHTTP", 404)
	}
	// }
}

func (p *Proxy) serveExtraRouteRequest(w http.ResponseWriter, req *http.Request) {
	defer common.TimeTrack(time.Now())
	var proxiedRequest bool
	routeAliases := p.aliasHosts[req.Host]
	if routeAliases != nil {
		apiServerName := p.hostsToAPIServer[req.Host]
		for _, v := range routeAliases {
			if strings.Contains(req.URL.Path, v) {
				p.proxyRequest(w, req, p.apiRoutes[apiServerName], req.URL.Path, req.Host, false)
				proxiedRequest = true
				break
			}
		}
		if proxiedRequest == false {
			p.invalidRoute(w, req.URL.String())
		}
	} else {
		p.invalidRoute(w, req.URL.String())
	}
}
func (p *Proxy) serveAPIRequest(w http.ResponseWriter, req *http.Request) {
	defer common.TimeTrack(time.Now())
	var apiName string
	var serviceAddr string
	var hostToReplace string
	if req.Host == p.apiBaseURL {
		apiName = p.getAPIName(req.URL.String(), hostToReplace)
		hostToReplace = p.apiBaseURLWithScheme
		serviceAddr = p.apiRoutes[apiName]
		if serviceAddr == "" {
			p.invalidRoute(w, req.URL.String())
			return
		}
	} else {
		p.invalidRoute(w, req.URL.String())
		return
	}
	p.proxyRequest(w, req, serviceAddr, "/api"+req.URL.Path, p.apiNameToOrigin[apiName], false)
}
func (p *Proxy) getAPIName(url, hostname string) string {
	var urlWithoutHost = strings.Replace(url, hostname, "", -1)
	var urlBits = strings.Split(urlWithoutHost, "/")
	return urlBits[1]
}
func (p *Proxy) proxyRequest(w http.ResponseWriter, req *http.Request, proxyTo string, path string, origin string, isRevProxReq bool) {
	var userInfo *url.Userinfo
	scheme := req.URL.Scheme
	if req.URL.User != nil {
		userInfo = req.URL.User
	}
	if forward.IsWebsocketRequest(req) == false {
		if isRevProxReq {
			scheme = "https"
		} else {
			scheme = "http"
		}
	} else {
		scheme = "ws"
	}
	req.URL = &url.URL{
		Host:     proxyTo,
		Path:     path,
		Scheme:   scheme,
		User:     userInfo,
		RawQuery: req.URL.Query().Encode(),
	}
	req.RequestURI = req.URL.String()
	if origin != "" {
		common.LogDebug("", "", origin)
		req.Header.Add("Access-Control-Allow-Origin", origin)
	}
	p.fwd.ServeHTTP(w, req)
}
func (p *Proxy) setRoutes() {
	var route string
	p.knownRoutes[p.apiBaseURL] = true
	p.knownRoutes[p.baseWDURL] = true
	p.knownRoutes[p.baseAuthURL] = true
	p.apiRoutes[watchdogAPIName] = "localhost:1000"
	p.apiRoutes[authServiceAPIName] = "localhost:1003"
	if routes, err := p.data.GetServiceDetailss(); err == nil {
		for _, v := range routes {
			if v.Internal {
				p.internal.SetInternalRoute(v.APIName, v.ServiceAddress)
			} else {
				p.apiRoutes[v.APIName] = v.ServiceAddress
				p.knownRoutes[v.AppName+p.baseURL] = true
				println(v.AppName + p.baseURL)
				p.apiNameToOrigin[v.APIName] = v.AppName + p.baseURL
			}
		}
	}

	if routes, e2 := p.data.GetAllExtraRoutes(); e2 == nil {
		for _, r := range routes {
			if strings.Contains(r.APIRoute, "*") {
				route = r.APIRoute[0 : len(r.APIRoute)-1]
			} else {
				route = r.APIRoute
			}
			p.knownAliasHosts[r.FullURL] = true
			p.hostsToAPIServer[r.FullURL] = r.APIName
			p.aliasHosts[r.FullURL] = append(p.aliasHosts[r.FullURL], route)
		}
	}

	// if revProxRoutes, e3 := p.data.GetProxyRoutes(); e3 == nil {
	// 	for _, r := range revProxRoutes {

	// 		p.reverseProxyRoutes[r.Hostname] = r.IPAddress
	// 	}
	// }
}
func (p *Proxy) invalidRoute(w http.ResponseWriter, requestedURL string) {
	//TOOD: Proper 404 page.
	common.WriteFailureResponse(errors.New("unknown route "+requestedURL), w, "ServeHTTP", 404)
}
func (p *Proxy) startTLSServer() {
	tlsConf := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	httpServer := &http.Server{
		Addr:         p.listenerPort,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		TLSConfig:    tlsConf,
		Handler:      p,
	}

	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	err = httpServer.ServeTLS(listener, "cert.pem", "cert.key")

	if err != nil {
		common.Logger.WithField("func", "main").Errorln(err)
	}
}
func (p *Proxy) startNotTLSServer() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	httpServer := &http.Server{
		Addr:         p.listenerPort,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      p,
	}
	err = httpServer.Serve(listener)
}
