package main

import (
	"os"

	"git.m/watchdog/common"
	"git.m/watchdog/data"
	"git.m/watchdog/management"
	"git.m/watchdog/proxy"
)

func main() {
	common.CommonProcessInit()
	common.Logger.Debugln("starting watchdog...")
	common.Logger.Debugln(os.Getuid())

	data := data.NewDataStoreInstance("routes")
	proxy := proxy.NewProxy(data)

	//TODO: Remove in favor of proper config UI.
	// data.AddNewRoute("gemini", "uihost", "gemini.rdro.us", "1998")
	// data.AddNewRoute("gemini-svc", "webservices", "gemini-svc.m.rdro.us", "1997")
	// data.AddNewRoute("player3", "p3uihost", "player3.m.rdro.us", "1999")

	manager := management.NewAPIRouter(data)
	manager.StartManagementAPIListener()

	proxy.SetRoutes()
	proxy.CreateListener()
}
