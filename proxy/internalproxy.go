package proxy

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vulcand/oxy/forward"
	"go.alargerobot.dev/frost/common"
)

//InternalProxy ...
type InternalProxy struct {
	listenerPort              string
	fwd                       *forward.Forwarder
	isInLocalMode             bool
	icAPIURL                  string
	internalBaseURL           string
	internalAPIBase           string
	internalAPIBaseWithScheme string
	internalRoutes            map[string]string
}

const (
	icapiAPIName = "icapi"
)

//NewInternalProxy ...
func NewInternalProxy() *InternalProxy {
	var fwd *forward.Forwarder
	sieh := &HTTPErrorHandler{}
	fwd, _ = forward.New(forward.ErrorHandler(sieh))
	return &InternalProxy{
		fwd: fwd,
		internalRoutes: make(map[string]string),
	}
}

//StartProxyListener ...
func (p *InternalProxy) StartProxyListener(localMode *bool) {
	common.LogInfo("", "", "starting internal service listener")
	p.internalRoutes[icapiAPIName] = "localhost:5000"
	p.internalBaseURL = "frost.m"
	if *localMode {
		common.Logger.Infoln("running in dev mode...")
		p.listenerPort = "8080"
		p.internalAPIBase = "api.frost-int.m"
	} else {
		p.listenerPort = "80"
		p.internalAPIBase = "api.frost.m"
	}
	p.startListener()
}

//SetInternalRoute ...
func (p *InternalProxy) SetInternalRoute(api, address string) {
	p.internalRoutes[api] = address
}

//ServeHTTP ...
func (p *InternalProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer common.TimeTrack(time.Now())
	var proxyTo string
	var userInfo *url.Userinfo

	scheme := req.URL.Scheme
	if req.URL.User != nil {
		userInfo = req.URL.User
	}
	if forward.IsWebsocketRequest(req) == false {
		scheme = "http"
	} else {
		scheme = "ws"
	}

	proxyTo = p.internalRoutes[p.getAPIName(req.URL.String(), p.internalAPIBase)]

	req.URL = &url.URL{
		Host:     proxyTo,
		Path:     "/api" + req.URL.Path,
		Scheme:   scheme,
		User:     userInfo,
		RawQuery: req.URL.Query().Encode(),
	}
	req.RequestURI = req.URL.String()

	p.fwd.ServeHTTP(w, req)
}

func (p *InternalProxy) getAPIName(url, hostname string) string {
	var urlWithoutHost = strings.Replace(url, hostname, "", -1)
	var urlBits = strings.Split(urlWithoutHost, "/")
	return urlBits[1]
}
func (p *InternalProxy) startListener() {
	listener, err := net.Listen("tcp", ":"+p.listenerPort)
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
