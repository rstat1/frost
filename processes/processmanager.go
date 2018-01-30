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
func (pm *ProcessManager) StartProcess(name, dirName string) bool {
	dir, _ := os.Getwd()
	path := dir + "/" + dirName + "/" + name
	common.Logger.Debugln(path)
	if _, err := os.Stat(dir + "/" + dirName + "/" + name); os.IsNotExist(err) {
		common.Logger.Debugln("not starting process...")
		return false
	} else {
		process := NewManagedProcess(path, []string{name, "-ppid", strconv.Itoa(os.Getpid())})
		pm.managedProcesses[name] = process
		process.Run()
		return true
	}
}

//StopAllProcesses ...
func (pm *ProcessManager) StopAllProcesses() {
	for _, v := range pm.managedProcesses {
		v.Stop <- true
	}
}

//StopAProcess ...
func (pm *ProcessManager) StopAProcess(name string) {
	process := pm.managedProcesses[name]
	process.Stop <- true
}
