package processes

import (
	"os"
	"strconv"

	"git.m/svcman/common"
)

//ProcessManager ...
type ProcessManager struct {
	managedProcesses map[string]*ManagedProcess
}

//NewProcessManager ...
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		managedProcesses: make(map[string]*ManagedProcess),
	}
}

//StartProcess ...
func (pm *ProcessManager) StartProcess(name, dirName string, devmode bool) bool {
	dir, _ := os.Getwd()
	path := dir + "/" + dirName + "/" + name
	if pm.managedProcesses[name] == nil {
		if _, err := os.Stat(dir + "/" + dirName + "/" + name); os.IsNotExist(err) {
			common.Logger.Debugln("not starting process...")
			return false
		} else {
			process := NewManagedProcess(path, dirName, []string{
				name, "-ppid", strconv.Itoa(os.Getpid()), "-devmode", strconv.FormatBool(devmode),
			})
			pm.managedProcesses[name] = process
			process.Run()
			return true
		}
	} else {
		return false
	}
}

//StopAllProcesses ...
func (pm *ProcessManager) StopAllProcesses() {
	for k, v := range pm.managedProcesses {
		v.Stop <- true
		delete(pm.managedProcesses, k)
	}
}

//StopAProcess ...
func (pm *ProcessManager) StopAProcess(name string) {
	process := pm.managedProcesses[name]
	if process != nil {
		process.Stop <- true
		delete(pm.managedProcesses, name)
	}
}