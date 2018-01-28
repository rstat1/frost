package processes

import (
	"os"
	"strconv"
)

//ProcessManager ...
type ProcessManager struct {
	managedProcesses map[string]ManagedProcess
}

//NewProcessManager ...
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		managedProcesses: make(map[string]ManagedProcess),
	}
}

//StartProcess ...
func (pm *ProcessManager) StartProcess(name string) {
	process := NewManagedProcess(name, []string{name, "-ppid", strconv.Itoa(os.Getpid())})
	pm.managedProcesses["name"] = *process
	process.Run()
}

//StopAllProcesses ...
func (pm *ProcessManager) StopAllProcesses() {
	for k := range pm.managedProcesses {
		pm.managedProcesses[k].Stop <- true
	}
}

//StopAProcess ...
func (pm *ProcessManager) StopAProcess(name string) {
	pm.managedProcesses[name].Stop <- true
}
