package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	//"strings"
	"time"

	"git.m/svcmanager/common"
	"git.m/svcmanager/data"
	"git.m/svcmanager/services"
	//"github.com/vulcand/oxy/utils"
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/vulcand/oxy/forward"
	"golang.org/x/crypto/acme/autocert"
)

//Proxy ...
type Proxy struct {
	data          *data.DataStore
	fwd           *forward.Forwarder
	knownUIRoutes []string
	knownRoutes   map[string]bool
	apiRoutes     map[string]string
	webuiserver   *services.WebServer
}

var (
	httpServer = &http.Server{
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	lastHost, lastHostPort string
)

const (
	baseURL              = ".dev-m.rdro.us"
	apiBaseURL           = "api" + baseURL
	watchdogURL          = "watchdog" + baseURL
	apiBaseURLWithScheme = "https://api" + baseURL

	watchdogAPIName = "frost"
	listenerPort    = ":80"
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
func (p *Proxy) StartProxyListener() {
	p.setRoutes()
	//m := autocert.Manager{
	//	Prompt:      autocert.AcceptTOS,
	//	Cache:       autocert.DirCache("certcache"),
	//	RenewBefore: 5 * time.Hour,
	//	HostPolicy:  p.urlWhiteList(),
	//	Email:       "rstat1@gmail.com",
	//}
	//go http.ListenAndServe(":http", m.HTTPHandler(nil))
	httpServer.TLSConfig = &tls.Config{
	//GetCertificate: m.GetCertificate,
	}
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	httpServer.Handler = p
	httpServer.Addr = listenerPort
	err = httpServer.Serve(listener)
	//err = httpServer.ServeTLS(listener, "", "")

	if err != nil {
		common.Logger.WithField("func", "main").Errorln(err)
	}
}

//AddRoute ...
func (p *Proxy) AddRoute(apiName, appName, serviceAddress string) {
	p.apiRoutes[apiName] = serviceAddress
	p.knownRoutes[appName+baseURL] = true
}

//ServeHTTP ...
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.String() == "/favicon.ico" {
		w.WriteHeader(200)
	} else if req.Host == apiBaseURL {
		p.serveAPIRequest(w, req)
	} else if strings.HasSuffix(req.Host, baseURL) {
		if p.knownRoutes[req.Host] {
			var name = strings.Replace(req.Host, baseURL, "", -1)
			p.webuiserver.ServeWebRequest(w, req, name)
		} else {
			//TOOD: Proper 404 page.
			w.WriteHeader(404)
		}
	} else {
		common.WriteFailureResponse(errors.New("not found"), w, "ServeHTTP", 404)
	}
}
func (p *Proxy) serveAPIRequest(w http.ResponseWriter, req *http.Request) {
	if req.Host == apiBaseURL {
		var urlWithoutHost = strings.Replace(req.URL.String(), apiBaseURLWithScheme, "", -1)
		var urlBits = strings.Split(urlWithoutHost, "/")
		var apiName = urlBits[1]
		req.URL = &url.URL{
			Host:     p.apiRoutes[apiName],
			Path:     "/api" + req.URL.Path,
			Scheme:   "http",
			RawQuery: req.URL.Query().Encode(),
		}
		req.RequestURI = req.URL.String()
		common.Logger.Debugln(req.URL.String())
		req.Header.Add("Access-Control-Allow-Origin", p.data.Cache.GetString("watchdog", apiName))
		p.fwd.ServeHTTP(w, req)
	}
}
func (p *Proxy) isWatchdogURL(url string) bool {
	if strings.HasPrefix(url, "svcmanager") || strings.Contains(url, watchdogAPIName) {
		return true
	}
	return false
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
	p.knownRoutes[apiBaseURL] = true
	p.knownRoutes[watchdogURL] = true
	p.apiRoutes[watchdogAPIName] = "localhost:1000"
	if routes, err := p.data.GetKnownRoutes(); err == nil {
		for _, v := range routes {
			p.apiRoutes[v.APIName] = v.ServiceAddress
			p.knownRoutes[v.AppName+baseURL] = true
			p.data.Cache.PutString("watchdog", v.APIName, v.AppName+baseURL)
		}
	}
}
