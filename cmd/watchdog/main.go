package main

import (
	"os"

	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/management"
	"git.m/svcman/proxy"
)

func main() {
	common.CommonProcessInit()
	common.Logger.Debugln("starting svcman...")
	common.Logger.Debugln(os.Getuid())

	data := data.NewDataStoreInstance("routes")
	proxy := proxy.NewProxy(data)
	services := management.NewServiceManager(data)
	manager := management.NewAPIRouter(data, proxy, services)

	services.StartManagedServices()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener()
}
