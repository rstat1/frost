package main

import (
	"flag"
	"os"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/management"
	"git.m/svcman/proxy"
)

func main() {
	common.CommonProcessInit()

	devMode := flag.Bool("devmode", false, "switches ports/URLs to dev mode, and also disables TLS support")
	flag.Parse()

	common.Logger.Debugln("starting svcman...")
	common.Logger.Debugln(os.Getuid())

	data := data.NewDataStoreInstance("routes")
	proxy := proxy.NewProxy(data)
	services := management.NewServiceManager(data, proxy)
	manager := management.NewAPIRouter(data, proxy, services)

	services.StartManagedServices()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener(*devMode)
}
