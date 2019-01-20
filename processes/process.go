package processes

import (
	"os"
	"syscall"

	"git.m/svcman/common"
)

type ManagedProcess struct {
	restartCount int
	WorkDir      string
	Name         string
	Args         []string
	restart      bool
	process      *os.Process
	Stop         chan bool
	Died         chan bool
	Stopped      bool
}

//NewManagedProcess ...
func NewManagedProcess(name, workDir string, args []string) *ManagedProcess {
	return &ManagedProcess{
		Name:    name,
		Args:    args,
		WorkDir: workDir,
	}
}

//Run ...
func (mp *ManagedProcess) Run() {
	var err error
	mp.restart = true
	procAttr := new(os.ProcAttr)
	procAttr.Dir = mp.WorkDir
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	procAttr.Env = []string{"PWD=" + mp.WorkDir}
	go func() {
	procloop:
		mp.Stop = make(chan bool)
		mp.Died = make(chan bool)
		mp.process, err = os.StartProcess(mp.Name, mp.Args, procAttr)
		if err != nil {
			common.Logger.Errorln(err)
			return
		}
		mp.restartCount++
		go func() {
			mp.process.Wait()
			mp.Died <- true
		}()
		select {
		case <-mp.Stop:
			mp.restart = false
			mp.Stopped = true
			mp.process.Signal(syscall.SIGTERM)
			return
		case <-mp.Died:
			if mp.restartCount <= 10 {
				if mp.restart == true {
					common.Logger.Debugln("restarting...")
					goto procloop
				}
			} else {
				common.Logger.Errorln("restart failed 10 times.")
				mp.Stopped = true
			}
		}
	}()
}
