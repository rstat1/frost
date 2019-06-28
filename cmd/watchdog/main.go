package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.m/svcman/auth"
	"git.m/svcman/common"
	"git.m/svcman/data"
	"git.m/svcman/management"
	"git.m/svcman/proxy"
)

func main() {
	devMode := flag.Bool("devmode", false, "switches ports/URLs to dev mode, and also disables TLS support")
	flag.Parse()

	common.CommonProcessInit(*devMode, true)

	common.Logger.Debugln("starting svcman...")
	common.Logger.Debugln(os.Getuid())

	data := data.NewDataStoreInstance("routes")
	userService := auth.NewUserService(data)

	proxy := proxy.NewProxy(data)
	services := management.NewServiceManager(data, proxy, *devMode)
	manager := management.NewAPIRouter(data, proxy, services, userService, *devMode)
	authService := auth.NewAuthService(data, userService, *devMode)

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		common.Logger.Debugln("exiting...")
		services.StopManagedServices()
		os.Exit(0)
	}()

	services.StartManagedServices()
	authService.InitAuthService()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener(devMode)
}

func interruptHandler() {
}
