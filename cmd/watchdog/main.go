package main

import (
	"flag"
	"os"

	"git.m/svcman/auth"
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
	userService := auth.NewUserService(data)

	proxy := proxy.NewProxy(data)
	services := management.NewServiceManager(data, proxy, *devMode)
	manager := management.NewAPIRouter(data, proxy, services, userService, *devMode)
	authService := auth.NewAuthService(data, userService, *devMode)

	services.StartManagedServices()
	authService.InitAuthService()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener(devMode)
}
