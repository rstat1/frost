package management

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"go.alargerobot.dev/frost/crypto"
	"go.alargerobot.dev/frost/processes"
	"go.alargerobot.dev/frost/proxy"

	"go.alargerobot.dev/frost/common"
	"go.alargerobot.dev/frost/data"
)

//ServiceManager ...
type ServiceManager struct {
	data      *data.DataStore
	proxy     *proxy.Proxy
	inDevMode bool
	vault     *crypto.VaultClient
	processes *processes.ProcessManager
}

//NewServiceManager ...
func NewServiceManager(store *data.DataStore, p *proxy.Proxy, devmode bool, vault *crypto.VaultClient) *ServiceManager {
	return &ServiceManager{
		proxy:     p,
		data:      store,
		vault:     vault,
		inDevMode: devmode,
		processes: processes.NewProcessManager(vault),
	}
}

//DeleteService ...
func (s *ServiceManager) DeleteService(name string) common.APIResponse {
	route, _ := s.data.GetRoute(name)
	s.proxy.DeleteRoute(route.APIName, route.AppName)
	s.data.DeleteRoute(name, false)

	return common.CreateAPIResponse("success", os.RemoveAll(route.AppName), 500)
}

//GetAllServices ...
func (s *ServiceManager) GetAllServices() []data.ServiceDetails {
	if routes, err := s.data.GetServiceDetailss(); err == nil {
		return routes
	} else {
		common.CreateFailureResponse(err, "GetAllServices", 500)
		return nil
	}
}

//GetServiceNames ...
func (s *ServiceManager) GetServiceNames() []string {
	var names []string
	if routes, err := s.data.GetServiceDetailss(); err == nil {
		for _, v := range routes {
			names = append(names, v.AppName)
		}
		return names
	} else {
		common.CreateFailureResponse(err, "GetAllServices", 500)
		return nil
	}
}

//StartManagedServices ...
func (s *ServiceManager) StartManagedServices() {
	go func() {
		if s.vault.TokenSet == false {
			common.LogInfo("", "", "waiting for Vault token to be set before continuning...")
			<-s.vault.TokenSetWatch
			common.LogInfo("", "", "Vault token set. Continuning with service start...")
		}
		for _, v := range s.GetAllServices() {
			if v.IsManagedService {
				s.StartManagedService(v.AppName)
			}
		}
	}()
}

//StopManagedServices ...
func (s *ServiceManager) StopManagedServices() {
	s.processes.StopAllProcesses()
}

//StartManagedService ...
func (s *ServiceManager) StartManagedService(name string) bool {
	routeInfo, _ := s.data.GetRoute(name)
	return s.processes.StartProcess(routeInfo.BinName, routeInfo.AppName, routeInfo.ServiceID, routeInfo.ServiceKey, s.inDevMode)
}

//StopManagedService ...
func (s *ServiceManager) StopManagedService(name string) {
	routeInfo, _ := s.data.GetRoute(name)
	s.processes.StopAProcess(routeInfo.BinName)
}

//NewService ...
func (s *ServiceManager) NewService(request *http.Request) common.APIResponse {
	resp, service := s.handleFileUpload(request, data.ServiceDetails{})
	if service.APIName != "" {
		s.StartManagedService(service.AppName)
	}
	return resp
}

//AddNewExtraRoute ...
func (s *ServiceManager) AddNewExtraRoute(newRoute data.RouteAlias) common.APIResponse {
	var resp common.APIResponse
	route := data.ExtraRoute{
		APIName:  newRoute.APIName,
		FullURL:  newRoute.FullURL,
		APIRoute: newRoute.APIRoute,
	}
	if e := s.data.AddExtraRoute(route); e != nil {
		resp = common.CreateAPIResponse("failed", e, 500)
	} else {
		resp = common.CreateAPIResponse("success", nil, 200)
		s.proxy.AddExtraRoute(route)
	}
	return resp
}

//UpdateService ...
func (s *ServiceManager) UpdateService(request *http.Request) common.APIResponse {
	serviceName := request.URL.Query().Get("name")
	if serviceName == "" {
		return common.CreateFailureResponse(errors.New("service name not specified"), "UpdateService", 500)
	}
	if serviceName == "watchdog" || serviceName == "trinity" {
		service := data.ServiceDetails{
			AppName: serviceName,
		}
		resp, _ := s.handleFileUpload(request, service)
		return resp
	}
	s.StopManagedService(serviceName)
	service, _ := s.data.GetRoute(serviceName)
	resp, _ := s.handleFileUpload(request, service)
	s.StartManagedService(serviceName)
	return resp
}

//GetService ...
func (s *ServiceManager) GetService(name string) (data.ServiceDetails, error) {
	if service, err := s.data.GetRoute(name); err == nil {
		return service, nil
	} else {
		return data.ServiceDetails{}, err
	}
}

//GetExtraRoutes ...
func (s *ServiceManager) GetExtraRoutes(apiName string) common.APIResponse {
	var resp common.APIResponse
	if routes, err := s.data.GetExtraRoutesForAPIName(apiName); err != nil {
		resp = common.CreateFailureResponse(err, "GetExtraRoutes", 400)
	} else {
		asJSON, _ := json.Marshal(routes)
		resp = common.CreateAPIResponse(string(asJSON), nil, 200)
	}
	return resp
}

//RenameServiceDirectory ...
func (s *ServiceManager) RenameServiceDirectory(oldname, newname string) error {
	if _, err := os.Stat(oldname); os.IsNotExist(err) {
		return err
	}
	s.StopManagedService(newname)
	e := os.Rename(oldname, newname)
	s.StartManagedService(newname)
	return e
}

//RenameServiceBin ...
func (s *ServiceManager) RenameServiceBin(oldName, newName, appName string) error {
	if _, err := os.Stat(appName + "/" + oldName); os.IsNotExist(err) {
		return err
	}
	s.processes.StopAProcess(oldName)
	e := os.Rename(appName+"/"+oldName, appName+"/"+newName)
	s.StartManagedService(newName)
	// s.processes.StartProcess(newName, appName, s.inDevMode)
	return e
}

func (s *ServiceManager) handleFileUpload(request *http.Request, info data.ServiceDetails) (common.APIResponse, data.ServiceDetails) {
	var err error
	var resp = common.CreateAPIResponse("success", nil, 500)
	var service data.ServiceDetails

	if err = request.ParseMultipartForm(75 * 1024 * 1024); err == nil {
		uiFiles, handler, noUIBlob := request.FormFile("uiblob")
		serviceFile, serviceFileInfo, notServiceBlob := request.FormFile("service")
		iconFile, _, noIcon := request.FormFile("icon")
		serviceDetails := request.FormValue("details")
		if serviceDetails != "" {
			json.Unmarshal([]byte(serviceDetails), &service)
			s.proxy.AddRoute(service)
			err = s.data.AddNewRoute(service)
			s.StartManagedService(service.AppName)
			if err != nil {
				resp = common.CreateFailureResponse(err, "NewService(AddNewRoute)", 500)
			}
		} else {
			service = info
		}
		if _, err := os.Stat(service.AppName); os.IsNotExist(err) {
			if err := os.Mkdir(service.AppName, 0700); err != nil {
				return common.CreateFailureResponse(err, "NewService(mkdir)", 500), data.ServiceDetails{}
			}
		}
		if notServiceBlob == nil {
			if strings.HasSuffix(serviceFileInfo.Filename, ".zip") {
				if err = s.handleServiceZipUpload(serviceFile, serviceFileInfo.Filename, service); err != nil {
					resp = common.CreateFailureResponse(err, "NewService(save-service)", 500)
				}
			} else {
				if err = s.handleServiceBinUpload(serviceFile, service.AppName+"/"+service.BinName); err != nil {
					resp = common.CreateFailureResponse(err, "NewService(save-service)", 500)
				}
			}
		}
		if noUIBlob == nil {
			if strings.HasSuffix(handler.Filename, ".zip") {
				if err = s.handleUIBlobUpload(uiFiles, handler.Filename, service.AppName); err != nil {
					resp = common.CreateFailureResponse(err, "NewService(save-service)", 500)
				}
			} else {
				resp = common.CreateFailureResponse(errors.New("uploaded ui blob not a zip"), "NewService(save-service)", 500)
			}
		}
		if noIcon == nil {
			if err = s.handleIconFileUpload(iconFile, service.AppName); err != nil {
				resp = common.CreateFailureResponse(err, "NewService(save-icon)", 500)
			}
		}
	} else {
		resp = common.CreateFailureResponse(err, "NewService(ParseForm)", 500)
	}
	return resp, service
}
func (s *ServiceManager) handleServiceBinUpload(fileContent multipart.File, fileName string) error {
	var serviceFileBytes bytes.Buffer
	if file, err := os.Create(fileName); err == nil {
		io.Copy(&serviceFileBytes, fileContent)
		if _, err := file.Write(serviceFileBytes.Bytes()); err != nil {
			common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
			return err
		}
		file.Chmod(0760)
		file.Close()
	} else {
		common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
		return err
	}
	return nil
}
func (s *ServiceManager) handleServiceZipUpload(fileContent multipart.File, fileName string, service data.ServiceDetails) error {
	var serviceFileBytes bytes.Buffer
	if file, err := os.Create(fileName); err == nil {
		io.Copy(&serviceFileBytes, fileContent)
		if _, err := file.Write(serviceFileBytes.Bytes()); err != nil {
			common.CreateFailureResponse(err, "handleServiceBinUpload(file-write)", 500)
			return err
		}
		file.Close()
		name, _ := os.Getwd()
		if err := common.Unzip(fileName, name+"/"+service.AppName+"/"); err != nil {
			common.CreateFailureResponse(err, "handleUIBlobUpload(unzip)", 500)
			common.Logger.WithField("func", "handleUIBlobUpload(unzip)").Errorln(err)
			return err
		}
	}
	return nil
}
func (s *ServiceManager) handleUIBlobUpload(fileContent multipart.File, fileName, appName string) error {
	var serviceFileBytes bytes.Buffer
	file, err := os.Create(fileName)
	if err == nil {
		io.Copy(&serviceFileBytes, fileContent)
		if _, err := file.Write(serviceFileBytes.Bytes()); err != nil {
			common.CreateFailureResponse(err, "handleServiceBinUpload(file-write)", 500)
			return err
		}
		file.Close()
		name, _ := os.Getwd()
		if err := common.Unzip(fileName, name+"/"+appName+"/web"); err != nil {
			common.CreateFailureResponse(err, "handleUIBlobUpload(unzip)", 500)
			common.Logger.WithField("func", "handleUIBlobUpload(unzip)").Errorln(err)
			return err
		}

		return os.Remove(fileName)

	}
	common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
	common.Logger.WithField("func", "handleUIBlobUpload(create)").Errorln(err)
	return err
}
func (s *ServiceManager) handleIconFileUpload(icon multipart.File, service string) error {
	var iconBytes bytes.Buffer
	if _, e := os.Stat("watchdog/serviceicons"); os.IsNotExist(e) {
		if err := os.Mkdir("watchdog/serviceicons", 0700); err != nil {
			return err
		}
	}
	if file, err := os.Create("watchdog/serviceicons/" + service + ".png"); err == nil {
		io.Copy(&iconBytes, icon)
		if _, err := file.Write(iconBytes.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
