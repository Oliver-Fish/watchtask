package execmanager

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

//Task is used to reference spawned processes and to initally spawn them
type Task struct {
	CmdPath    string
	Args       []string
	Running    bool
	Cmd        *exec.Cmd
	StdoutPipe io.Reader
}

//Start takes a task struct and starts the command in a none blocking way and error if it fails
func (t *Task) Start() error {
	t.Cmd = exec.Command(t.CmdPath, t.Args...)
	r, err := t.Cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	}
	t.StdoutPipe = r
	t.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} //Sets a gobal pid for child processes credit to https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	err = t.Cmd.Start()
	if err != nil {
		return err
	}

	t.Running = true

	return nil
}

//Kill will attempt kill the gpid of the passed task and return and error if it fails
func (t *Task) Kill() error {
	//Todo check if pid is still valid https://stackoverflow.com/questions/15204162/check-if-a-process-exists-in-go-way
	err := syscall.Kill(-t.Cmd.Process.Pid, syscall.SIGKILL) //Kills global pid
	t.Running = false
	t.StdoutPipe = nil
	return err
}
