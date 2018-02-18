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
	knownUIRoutes        []string
	knownRoutes          map[string]bool
	apiRoutes            map[string]string
	webuiserver          *services.WebServer
	isInLocalMode        bool
	apiBaseURL           string
	baseURL              string
	apiBaseURLWithScheme string
	listenerPort         string
}

var (
	httpServer = &http.Server{
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	lastHost, lastHostPort string
)

const (
	watchdogAPIName = "frost"

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
		fwd:         fwd,
		data:        dataStoreRef,
		apiRoutes:   make(map[string]string),
		knownRoutes: make(map[string]bool),
		webuiserver: services.NewWebServer(),
	}
}

//StartProxyListener ...
func (p *Proxy) StartProxyListener(localMode bool) {
	p.isInLocalMode = localMode
	if localMode == false {
		p.baseURL = prodBaseURL
		p.apiBaseURL = prodBaseAPIURL
		p.listenerPort = productionPort
		p.apiBaseURLWithScheme = baseAPIURLWithScheme

		common.Logger.Infoln("running in production mode...")

		p.setRoutes()
		p.startTLSServer()
	} else {
		p.baseURL = devBaseURL
		p.apiBaseURL = devAPIBaseURL
		p.listenerPort = devPort
		p.apiBaseURLWithScheme = devAPIBaseURLWithScheme

		common.Logger.Infoln("running in dev mode...")

		p.setRoutes()
		p.startNotTLSServer()
	}
}

//AddRoute ...
func (p *Proxy) AddRoute(newRoute data.KnownRoute) {
	p.apiRoutes[newRoute.APIName] = newRoute.ServiceAddress
	p.knownRoutes[newRoute.AppName+p.baseURL] = true
}

//DeleteRoute ...
func (p *Proxy) DeleteRoute(apiName, appName string) {
	delete(p.apiRoutes, apiName)
	delete(p.knownRoutes, appName+p.baseURL)
}

//ServeHTTP ...
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.String() == "/favicon.ico" {
		w.WriteHeader(200)
	} else if req.Host == p.apiBaseURL {
		p.serveAPIRequest(w, req)
	} else if strings.HasSuffix(req.Host, p.baseURL) {
		if p.knownRoutes[req.Host] {
			var name = strings.Replace(req.Host, p.baseURL, "", -1)
			p.webuiserver.ServeWebRequest(w, req, name)
		} else {
			p.invalidRoute(w, req.Host)
		}
	} else {
		common.WriteFailureResponse(errors.New("not found"), w, "ServeHTTP", 404)
	}
}
func (p *Proxy) serveAPIRequest(w http.ResponseWriter, req *http.Request) {
	if req.Host == p.apiBaseURL {
		var urlWithoutHost = strings.Replace(req.URL.String(), p.apiBaseURLWithScheme, "", -1)
		var urlBits = strings.Split(urlWithoutHost, "/")
		var apiName = urlBits[1]
		if p.apiRoutes[apiName] != "" {
			scheme := req.URL.Scheme
			if forward.IsWebsocketRequest(req) == false {
				scheme = "http"
			} else {
				scheme = "ws"
				common.Logger.Debugln(req)
			}
			req.URL = &url.URL{
				Host:     p.apiRoutes[apiName],
				Path:     "/api" + req.URL.Path,
				Scheme:   scheme,
				RawQuery: req.URL.Query().Encode(),
			}
			req.RequestURI = req.URL.String()
			req.Header.Add("Access-Control-Allow-Origin", p.data.Cache.GetString("watchdog", apiName))
			p.fwd.ServeHTTP(w, req)

		} else {
			p.invalidRoute(w, req.Host)
		}
	}
}
func (p *Proxy) urlWhiteList() autocert.HostPolicy {
	return func(_ context.Context, host string) error {
		common.Logger.Debugln(host)
		if !p.knownRoutes[host] {
			err := errors.New("host not on whitelist: " + host)
			common.CreateFailureResponse(err, "urlWhiteList", 400)
			return err
		}
		return nil
	}
}
func (p *Proxy) setRoutes() {
	p.knownRoutes[p.apiBaseURL] = true
	p.knownRoutes[watchdogURL] = true
	p.apiRoutes[watchdogAPIName] = "localhost:1000"
	if routes, err := p.data.GetKnownRoutes(); err == nil {
		for _, v := range routes {
			p.apiRoutes[v.APIName] = v.ServiceAddress
			p.knownRoutes[v.AppName+p.baseURL] = true
			p.data.Cache.PutString("watchdog", v.APIName, v.AppName+p.baseURL)
		}
		common.Logger.Debugln(p.knownRoutes)
	}
}
func (p *Proxy) invalidRoute(resp http.ResponseWriter, requestedURL string) {
	//TOOD: Proper 404 page.
	resp.Write([]byte(requestedURL + " not found"))
	resp.WriteHeader(404)
}
func (p *Proxy) startTLSServer() {
	m := autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache("certcache"),
		RenewBefore: 5 * time.Hour,
		HostPolicy:  p.urlWhiteList(),
		Email:       "rstat1@gmail.com",
	}

	s := &http.Server{
		Handler: m.HTTPHandler(nil),
		Addr:    ":80",
	}
	go s.ListenAndServe()
	httpServer.TLSConfig = &tls.Config{
		GetCertificate: m.GetCertificate,
	}
	listener, err := net.Listen("tcp", p.listenerPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	httpServer.Handler = p
	httpServer.Addr = p.listenerPort
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
	httpServer.Handler = p
	httpServer.Addr = p.listenerPort
	err = httpServer.Serve(listener)
}
