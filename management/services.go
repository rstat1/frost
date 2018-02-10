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

	"git.m/svcman/processes"
	"git.m/svcman/proxy"

	"git.m/svcman/common"
	"git.m/svcman/data"
)

//ServiceManager ...
type ServiceManager struct {
	data      *data.DataStore
	proxy     *proxy.Proxy
	inDevMode bool
	processes *processes.ProcessManager
}

//NewServiceManager ...
func NewServiceManager(store *data.DataStore, p *proxy.Proxy, devmode bool) *ServiceManager {
	return &ServiceManager{
		proxy:     p,
		data:      store,
		inDevMode: devmode,
		processes: processes.NewProcessManager(),
	}
}

//DeleteService ...
func (s *ServiceManager) DeleteService(name string) common.APIResponse {
	route := s.data.GetRoute(name)
	s.proxy.DeleteRoute(route.APIName, route.AppName)
	s.data.DeleteRoute(name)
	return common.CreateAPIResponse("success", os.RemoveAll(route.AppName), 500)
}

//GetAllServices ...
func (s *ServiceManager) GetAllServices() []data.KnownRoute {
	if routes, err := s.data.GetKnownRoutes(); err == nil {
		return routes
	} else {
		common.CreateFailureResponse(err, "GetAllServices", 500)
		return nil
	}
}

//StartManagedServices ...
func (s *ServiceManager) StartManagedServices() {
	for _, v := range s.GetAllServices() {
		if v.IsManagedService {
			s.processes.StartProcess(v.BinName, v.AppName, s.inDevMode)
		}
	}
}

//StartManagedService ...
func (s *ServiceManager) StartManagedService(name string) bool {
	var routeInfo data.KnownRoute = s.data.GetRoute(name)
	return s.processes.StartProcess(routeInfo.BinName, routeInfo.AppName, s.inDevMode)
}

//StopManagedService ...
func (s *ServiceManager) StopManagedService(name string) {
	var routeInfo data.KnownRoute = s.data.GetRoute(name)
	s.processes.StopAProcess(routeInfo.BinName)
}

//NewService ...
func (s *ServiceManager) NewService(request *http.Request) common.APIResponse {
	resp, service := s.handleFileUpload(request, data.KnownRoute{})
	if service.APIName != "" {
		s.StartManagedService(service.AppName)
	}
	return resp
}

//UpdateService ...
func (s *ServiceManager) UpdateService(request *http.Request) common.APIResponse {
	serviceName := request.URL.Query().Get("name")
	if serviceName == "" {
		return common.CreateFailureResponse(errors.New("service name not specified"), "UpdateService", 500)
	} else {
		s.StopManagedService(serviceName)
		resp, _ := s.handleFileUpload(request, s.data.GetRoute(serviceName))
		s.StartManagedService(serviceName)
		return resp
	}
}
func (s *ServiceManager) handleFileUpload(request *http.Request, info data.KnownRoute) (common.APIResponse, data.KnownRoute) {
	var err error
	var resp common.APIResponse = common.CreateAPIResponse("success", nil, 500)
	var service data.KnownRoute

	if err = request.ParseMultipartForm(75 * 1024 * 1024); err == nil {
		uiFiles, handler, noUIBlob := request.FormFile("uiblob")
		serviceFile, _, notServiceBlob := request.FormFile("service")
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
			if err := os.Mkdir(service.AppName, 0770); err != nil {
				return common.CreateFailureResponse(err, "NewService(mkdir)", 500), data.KnownRoute{}
			}
		}
		if notServiceBlob == nil {
			if err = s.handleServiceBinUpload(serviceFile, service.AppName+"/"+service.BinName); err != nil {
				resp = common.CreateFailureResponse(err, "NewService(save-service)", 500)
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
func (s *ServiceManager) handleUIBlobUpload(fileContent multipart.File, fileName, appName string) error {
	var serviceFileBytes bytes.Buffer
	if file, err := os.Create(fileName); err == nil {
		io.Copy(&serviceFileBytes, fileContent)
		if _, err := file.Write(serviceFileBytes.Bytes()); err != nil {
			common.CreateFailureResponse(err, "handleServiceBinUpload(file-write)", 500)
			return err
		}
		file.Close()
		name, _ := os.Getwd()
		if err := common.Unzip(fileName, name+"/"+appName+"/web"); err != nil {
			common.CreateFailureResponse(err, "handleUIBlobUpload(unzip)", 500)
			return err
		} else {
			os.Remove(fileName)
			return nil
		}
	} else {
		common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
		return err
	}
}
