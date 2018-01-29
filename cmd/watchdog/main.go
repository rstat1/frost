package main

import (
	"os"

	"git.m/svcmanager/common"
	"git.m/svcmanager/data"
	"git.m/svcmanager/management"
	"git.m/svcmanager/proxy"
)

func main() {
	common.CommonProcessInit()
	common.Logger.Debugln("starting svcmanager...")
	common.Logger.Debugln(os.Getuid())

	data := data.NewDataStoreInstance("routes")
	proxy := proxy.NewProxy(data)
	services := management.NewServiceManager(data)
	manager := management.NewAPIRouter(data, proxy, services)

	services.StartManagedServices()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener()
}
