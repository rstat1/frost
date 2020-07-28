package processes

import (
	"os"
	"strconv"
	"syscall"

	"go.alargerobot.dev/frost/common"
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
func (pm *ProcessManager) StartProcess(name, dirName, sid, skey string, devmode bool) bool {
	dir, _ := os.Getwd()
	path := dir + "/" + dirName + "/" + name
	if pm.managedProcesses[name] == nil {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			common.Logger.WithField("path", path).Debugln("not starting process...")
			return false
		} else {
			process := NewManagedProcess(path, dirName, []string{
				name, "-ppid", strconv.Itoa(os.Getpid()), "-devmode=" + strconv.FormatBool(devmode),
			})
			pm.managedProcesses[name] = process
			process.Run(sid, skey)
			return true
		}
	} else {
		return false
	}
}

//StopAllProcesses ...
func (pm *ProcessManager) StopAllProcesses() {
	for _, v := range pm.managedProcesses {
		v.process.Signal(syscall.SIGTERM)
	}
}

//StopAProcess ...
func (pm *ProcessManager) StopAProcess(name string) {
	process := pm.managedProcesses[name]
	if process != nil {
		if process.Stopped == false {
			process.Stop <- true
		}
		delete(pm.managedProcesses, name)
	}
}
