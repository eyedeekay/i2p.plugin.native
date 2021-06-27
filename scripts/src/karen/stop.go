package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func Stop() error {
	errstr := ""
	if cmd != nil {
		Printf("Stopping %d %s %s", cmd.Process.Pid, *executableFile, getArgs())
		if runtime.GOOS == "windows" {
			err := cmd.Process.Signal(os.Kill)
			if err != nil {
				errstr += err.Error()
			}
		} else {
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				errstr += err.Error()
			}
		}
		if errstr == "" {
			os.Remove(filepath.Join(*runDir, strconv.Itoa(cmd.Process.Pid)))
		}
	}
	err := StopCommands()
	if err != nil {
		errstr += err.Error()
	}
	if errstr == "" {
		return nil
	}
	return fmt.Errorf("%s", errstr)
}
