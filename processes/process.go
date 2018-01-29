package processes

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"git.m/svcmanager/common"
)

type ManagedProcess struct {
	Name         string
	Args         []string
	Stop         chan bool
	restart      bool
	process      *os.Process
	restartCount int
}

//NewManagedProcess ...
func NewManagedProcess(name string, args []string) *ManagedProcess {
	return &ManagedProcess{
		Name: name,
		Args: args,
		Stop: make(chan bool, 1),
	}
}

//Run ...
func (mp *ManagedProcess) Run() {
	var err error
	var ws syscall.WaitStatus

	childProcessSignal := make(chan os.Signal, 1)
	//signal.Notify(childProcessSignal, syscall.SIGCHLD)
	mp.restart = true

	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	go func() {
		common.Logger.WithField("processname", mp.Name).Debugln("started process")
		defer func() {
			signal.Stop(childProcessSignal)
		}()
	procloop:
		for {
			if mp.restart {
				if mp.restartCount <= 10 {
					common.SentryClient.CapturePanic(func() {
						if mp.process, err = os.StartProcess(mp.Name, mp.Args, procAttr); err != nil {
							common.Logger.Errorln(err)
							return
						}
						mp.restart = false
						mp.restartCount++
					}, nil, nil)
				} else {
					common.Logger.WithField("processname", mp.Name).Errorln(errors.New("restart failed 10 times"))
					return
				}
				mp.restart = false
			}
			select {
			case <-childProcessSignal:
				common.Logger.Debugln("cps")
				pid, err := syscall.Wait4(mp.process.Pid, &ws, syscall.WNOHANG, nil)
				common.Logger.Debugln("cps2")
				if err != nil {
					common.Logger.Errorln(err)
				}
				if pid == mp.process.Pid {
					mp.restart = true
				}
				continue procloop
			case <-mp.Stop:
				mp.process.Kill()
				mp.restart = false
				return
			}
		}
	}()
}
