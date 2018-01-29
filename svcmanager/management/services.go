package management

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"git.m/svcmanager/processes"

	"git.m/svcmanager/common"
	"git.m/svcmanager/data"
)

//ServiceManager ...
type ServiceManager struct {
	data      *data.DataStore
	processes *processes.ProcessManager
}

//NewServiceManager ...
func NewServiceManager(store *data.DataStore) *ServiceManager {
	return &ServiceManager{
		data:      store,
		processes: processes.NewProcessManager(),
	}
}

//DeleteService ...
func (s *ServiceManager) DeleteService(name string) common.APIResponse {
	return common.CreateAPIResponse("success", s.data.DeleteRoute(name), 500)
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
			s.processes.StartProcess(v.BinName, v.AppName)
		}
	}
}

//StartManagedService ...
func (s *ServiceManager) StartManagedService(name string) bool {
	var routeInfo data.KnownRoute = s.data.GetRoute(name)
	return s.processes.StartProcess(routeInfo.BinName, routeInfo.AppName)
}

//StopManagedService ...
func (s *ServiceManager) StopManagedService(name string) {
	var routeInfo data.KnownRoute = s.data.GetRoute(name)
	s.processes.StopAProcess(routeInfo.BinName)
}

//NewService ...
func (s *ServiceManager) NewService(request *http.Request) common.APIResponse {
	var err error
	var service data.KnownRoute
	var resp common.APIResponse = common.CreateAPIResponse("success", nil, 500)

	serviceDetails := request.FormValue("details")

	json.Unmarshal([]byte(serviceDetails), &service)
	err = s.data.AddNewRoute(service)
	if err != nil {
		resp = common.CreateFailureResponse(err, "NewService(AddNewRoute)", 500)
	} else {
		//resp = s.handleFileUpload(request, "")
	}
	return resp
}
func (s *ServiceManager) UpdateService(request *http.Request) common.APIResponse {
	return common.APIResponse{}
}
func (s *ServiceManager) handleFileUpload(request *http.Request, uploadType, appName, binName string) error {
	var err error
	var serviceFileBytes bytes.Buffer

	if err = request.ParseMultipartForm(75 * 1024 * 1024); err == nil {
		serviceFile, _, _ := request.FormFile("service")
		io.Copy(&serviceFileBytes, serviceFile)
		if _, err := os.Stat(appName); os.IsNotExist(err) {
			if err := os.Mkdir(appName, 0644); err == nil {
				if uploadType == "service" {
					err = s.handleServiceBinUpload(serviceFileBytes.Bytes(), appName+"/"+binName)
				} else if uploadType == "ui" {
					err = s.handleUIBlobUpload(serviceFileBytes.Bytes(), "")
				}
			} else {
				common.CreateFailureResponse(err, "NewService(mkdir)", 500)
				return err
			}
		}
	} else {
		common.CreateFailureResponse(err, "NewService(ParseForm)", 500)
		return err
	}
	return nil
}
func (s *ServiceManager) handleServiceBinUpload(fileContent []byte, fileName string) error {
	if file, err := os.Create(fileName); err == nil {
		if _, err := file.Write(fileContent); err != nil {
			common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
			return err
		}
		file.Chmod(0744)
		file.Close()
	} else {
		common.CreateFailureResponse(err, "handleServiceBinUpload", 500)
		return err
	}
	return nil
}
func (s *ServiceManager) handleUIBlobUpload(fileContent []byte, name string) error {

	return nil
}
