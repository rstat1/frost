package processes

import (
	"os"
	"strconv"
	"syscall"

	"go.alargerobot.dev/frost/common"
	"go.alargerobot.dev/frost/crypto"
)

//ProcessManager ...
type ProcessManager struct {
	vault            *crypto.VaultClient
	managedProcesses map[string]*ManagedProcess
}

//NewProcessManager ...
func NewProcessManager(vc *crypto.VaultClient) *ProcessManager {
	return &ProcessManager{
		vault:            vc,
		managedProcesses: make(map[string]*ManagedProcess),
	}
}

//StartProcess ...
func (pm *ProcessManager) StartProcess(name, dirName, sid, skey string, devmode, useVault bool) bool {
	var envVars []string
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

			if useVault {
				id, err := pm.vault.GetRoleID(name)
				if err != nil {
					common.LogError("", err)
					return false
				}

				arsid, err := pm.vault.GetSecretIDAccessor()
				if err != nil {
					common.LogError("", err)
					return false
				}
				envVars = []string{"SKEY=" + skey, "SID=" + sid, "APPROLE_ID=" + id, "ARSID_ACCESS_KEY=" + arsid, "VAULTADDR=" + common.CurrentConfig.VaultAddr}
			} else {
				envVars = []string{"SKEY=" + skey, "SID=" + sid}
			}

			process.Run(envVars)
			return true
		}
	} else {
		return false
	}
}

//StopAllProcesses ...
func (pm *ProcessManager) StopAllProcesses() {
	for _, v := range pm.managedProcesses {
		if v.process != nil {
			v.process.Signal(syscall.SIGTERM)
		}
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
