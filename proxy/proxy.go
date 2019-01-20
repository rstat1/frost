package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"context"
	"net/url"
	"strings"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/services"

	"github.com/pkg/errors"
	"github.com/vulcand/oxy/forward"
	"golang.org/x/crypto/acme/autocert"
)

//Proxy ...
type Proxy struct {
	data                 *data.DataStore
	fwd                  *forward.Forwarder
	knownRoutes          map[string]bool
	apiRoutes            map[string]string
	aliasHosts           map[string][]string
	knownAliasHosts      map[string]bool
	hostsToAPIServer     map[string]string
	webuiserver          *services.WebServer
	isInLocalMode        bool
	apiBaseURL           string
	baseURL              string
	baseAuthURL          string
	baseWDURL            string
	apiBaseURLWithScheme string
	listenerPort         string
}

var (
	lastHost, lastHostPort string
)

const (
	watchdogAPIName    = "frost"
	authServiceAPIName = "trinity"

	prodBaseURL          = ".m.rdro.us"
	prodBaseAPIURL       = "api" + prodBaseURL
	watchdogURL          = "watchdog" + prodBaseURL
	baseAPIURLWithScheme = "https://" + prodBaseURL

	devBaseURL              = ".dev-m.rdro.us"
	devAPIBaseURL           = "api" + devBaseURL
	devWatchdogURL          = "watchdog" + devBaseURL
	devAPIBaseURLWithScheme = "http://" + devAPIBaseURL

	devPort        = ":80"
	productionPort = ":443"
)

//NewProxy ...
func NewProxy(dataStoreRef *data.DataStore) *Proxy {
	var fwd *forward.Forwarder
	fwd, _ = forward.New() //forward.Logger(common.Logger))
	return &Proxy{
		fwd:              fwd,
		data:             dataStoreRef,
		apiRoutes:        make(map[string]string),
		aliasHosts:       make(map[string][]string),
		knownRoutes:      make(map[string]bool),
		knownAliasHosts:  make(map[string]bool),
		webuiserver:      services.NewWebServer(),
		hostsToAPIServer: make(map[string]string),
	}
}

//StartProxyListener ...
func (p *Proxy) StartProxyListener(localMode *bool) {
	common.Logger.WithField("localmode", *localMode).Infoln("starting proxy listener.")

	p.isInLocalMode = *localMode
	if p.isInLocalMode == false {
		p.baseURL = prodBaseURL
		p.apiBaseURL = prodBaseAPIURL
		p.listenerPort = productionPort
		p.apiBaseURLWithScheme = baseAPIURLWithScheme
		p.baseAuthURL = "trinity" + prodBaseURL
		p.baseWDURL = watchdogURL
		common.Logger.Infoln("running in production mode...")

		p.setRoutes()
		p.startTLSServer()
	} else {
		p.baseURL = devBaseURL
		p.apiBaseURL = devAPIBaseURL
		p.listenerPort = devPort
		p.apiBaseURLWithScheme = devAPIBaseURLWithScheme
		p.baseAuthURL = "trinity" + devBaseURL
		p.baseWDURL = devWatchdogURL
		common.Logger.Infoln("running in dev mode...")

		p.setRoutes()
		p.startNotTLSServer()
	}
}

//AddRoute ...
func (p *Proxy) AddRoute(newRoute data.ServiceDetails) {
	p.apiRoutes[newRoute.APIName] = newRoute.ServiceAddress
	p.knownRoutes[newRoute.AppName+p.baseURL] = true
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
}
func (p *Proxy) serveExtraRouteRequest(w http.ResponseWriter, req *http.Request) {
	var proxiedRequest bool
	routeAliases := p.aliasHosts[req.Host]
	if routeAliases != nil {
		apiServerName := p.hostsToAPIServer[req.Host]
		for _, v := range routeAliases {
			// common.Logger.WithFields(logrus.Fields{"routeAlias": v, "reqPath": req.URL.Path}).Debugln("routeAliases")
			if strings.Contains(req.URL.Path, v) {
				p.proxyRequest(w, req, p.apiRoutes[apiServerName], req.URL.Path, req.Host)
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
	if req.Host == p.apiBaseURL {
		var urlWithoutHost = strings.Replace(req.URL.String(), p.apiBaseURLWithScheme, "", -1)
		var urlBits = strings.Split(urlWithoutHost, "/")
		var apiName = urlBits[1]
		p.proxyRequest(w, req, p.apiRoutes[apiName], "/api"+req.URL.Path, p.data.Cache.GetString("watchdog", apiName))
	} else {
		p.invalidRoute(w, req.URL.String())
	}
}
func (p *Proxy) proxyRequest(w http.ResponseWriter, req *http.Request, proxyTo string, path string, origin string) {
	scheme := req.URL.Scheme
	if forward.IsWebsocketRequest(req) == false {
		scheme = "http"
	} else {
		scheme = "ws"
	}
	req.URL = &url.URL{
		Host:     proxyTo,
		Path:     path,
		Scheme:   scheme,
		RawQuery: req.URL.Query().Encode(),
	}
	req.RequestURI = req.URL.String()
	if origin != "" {
		req.Header.Add("Access-Control-Allow-Origin", origin)
	}
	p.fwd.ServeHTTP(w, req)
}
func (p *Proxy) urlWhiteList() autocert.HostPolicy {
	return func(_ context.Context, host string) error {
		if !p.knownRoutes[host] && !p.knownAliasHosts[host] {
			err := errors.New("host not on whitelist: " + host)
			common.CreateFailureResponse(err, "urlWhiteList", 400)
			return err
		}
		return nil
	}
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
			p.apiRoutes[v.APIName] = v.ServiceAddress
			p.knownRoutes[v.AppName+p.baseURL] = true
			p.data.Cache.PutString("watchdog", v.APIName, v.AppName+p.baseURL)
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
}
func (p *Proxy) invalidRoute(w http.ResponseWriter, requestedURL string) {
	//TOOD: Proper 404 page.
	common.WriteFailureResponse(errors.New("unknown route "+requestedURL), w, "ServeHTTP", 404)
}
func (p *Proxy) startTLSServer() {
	m := autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache("certcache"),
		RenewBefore: 5 * time.Hour,
		HostPolicy:  p.urlWhiteList(),
		Email:       "rstat1@gmail.com",
	}

	// s := &http.Server{
	// 	Handler: m.HTTPHandler(nil),
	// 	Addr:    ":80",
	// }
	// go s.ListenAndServe()
	tlsConf := m.TLSConfig()
	tlsConf.MinVersion = tls.VersionTLS10
	httpServer := &http.Server{
		Addr:         p.listenerPort,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		TLSConfig:    tlsConf,
		Handler:      p,
	}

	listener, err := net.Listen("tcp", p.listenerPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	err = httpServer.ServeTLS(listener, "", "")

	if err != nil {
		common.Logger.WithField("func", "main").Errorln(err)
	}
}
func (p *Proxy) startNotTLSServer() {
	listener, err := net.Listen("tcp", p.listenerPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	httpServer := &http.Server{
		Addr:         p.listenerPort,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      p,
	}
	err = httpServer.Serve(listener)
}
