package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.alargerobot.dev/frost/auth"
	"go.alargerobot.dev/frost/common"
	"go.alargerobot.dev/frost/crypto"
	"go.alargerobot.dev/frost/data"
	"go.alargerobot.dev/frost/management"
	"go.alargerobot.dev/frost/proxy"
)

func main() {
	devMode := flag.Bool("devmode", false, "switches ports/URLs to dev mode, and also disables TLS support")
	flag.Parse()

	common.CommonProcessInit(*devMode, true)

	common.Logger.Debugln("starting svcman...")

	data := data.NewDataStoreInstance("routes")
	userService := auth.NewUserService(data)
	vault := crypto.NewVaultClient(*devMode)

	proxy := proxy.NewProxy(data)
	icapi := management.NewInternalConfigAPI(data, vault)
	services := management.NewServiceManager(data, proxy, *devMode, vault)
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

	icapi.InitListener()
	services.StartManagedServices()
	authService.InitAuthService()
	manager.StartManagementAPIListener()
	proxy.StartProxyListener(devMode)
}

func interruptHandler() {
}
