package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rstat1/frost/auth"
	"github.com/rstat1/frost/common"
	"github.com/rstat1/frost/crypto"
	"github.com/rstat1/frost/data"
	"github.com/rstat1/frost/management"
	"github.com/rstat1/frost/proxy"
)

func main() {
	devMode := flag.Bool("devmode", false, "switches ports/URLs to dev mode, and also disables TLS support")
	flag.Parse()

	common.CommonProcessInit(*devMode, true)

	common.Logger.Debugln("starting svcman...")

	data := data.NewDataStoreInstance("routes")
	userService := auth.NewUserService(data)	
	vault := crypto.NewVaultClient(*devMode)

	proxy := proxy.NewProxy(data, devMode)
	icapi := management.NewInternalPlatformAPI(data, vault, userService)
	services := management.NewServiceManager(data, proxy, *devMode, vault)
	manager := management.NewAPIRouter(data, proxy, services, userService, vault, *devMode)
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
	proxy.StartProxyListener()
}

func interruptHandler() {
}
