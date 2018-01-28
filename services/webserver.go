package services

import (
	"net/http"
	"strings"

	"git.m/watchdog/data"
)

//WebServer ...
type WebServer struct {
	data *data.DataStore
}

//NewWebServer ...
func NewWebServer() *WebServer {
	return &WebServer{}
}

//ServeWebRequest ...
func (web *WebServer) ServeWebRequest(w http.ResponseWriter, r *http.Request, appName string) {
	var pathToServe string
	if web.isFilePath(r.URL.Path) {
		pathToServe = appName + "/web/" + r.URL.Path
	} else {
		pathToServe = appName + "/web/index.html"
	}
	http.ServeFile(w, r, pathToServe)

}
func (web *WebServer) isFilePath(path string) bool {
	if strings.Contains(path, ".css") ||
		strings.Contains(path, ".js") ||
		strings.Contains(path, ".png") ||
		strings.Contains(path, ".jpg") {
		return true
	}
	return false
}
