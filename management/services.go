package management

import (
	"git.m/watchdog/common"
	"git.m/watchdog/data"
)

type ServiceManager struct {
	data *data.DataStore
}

//NewServiceManager ...
func NewServiceManager(store *data.DataStore) *ServiceManager {
	return &ServiceManager{
		data: store,
	}
}

//AddService ...
func (s *ServiceManager) AddService(service data.KnownRoute) common.APIResponse {
	return common.CreateAPIResponse("success", s.data.AddNewRoute(service), 500)
}
