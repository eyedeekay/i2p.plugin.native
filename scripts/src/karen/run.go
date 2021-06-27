package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var cmd *exec.Cmd

func Run() error {
	Printf("Starting %s %s", *executableFile, getArgs())
	cmd = exec.Command(getExecutable(), getArgs()...)
	cmd.Dir = *executableDir
	logfile, err := os.Create(filepath.Join(*executableDir, *commandLog))
	if err != nil {
		return err
	}
	cmd.Stdout = logfile
	errfile, err := os.Create(filepath.Join(*executableDir, *commandErr))
	if err != nil {
		return err
	}
	cmd.Stderr = errfile
	return cmd.Start()
}

func WritePID() error {
	spid := strconv.Itoa(cmd.Process.Pid)
	pid := []byte(strconv.Itoa(cmd.Process.Pid))
	err := ioutil.WriteFile(filepath.Join(*runDir, spid), pid, 0644)
	if err != nil {
		return err
	}
	return nil
}
