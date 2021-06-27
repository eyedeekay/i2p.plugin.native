package main

import (
	"fmt"
	"syscall"
)

func StatusMessage() error {
	err := GatherCommands()
	if err != nil {
		return err
	}
	if len(GatheredCommands) > 2 {
		return fmt.Errorf("Running with dangling PID's, please restart to fix\n if this error persists check your rundir permissions")
	}
	for _, command := range GatheredCommands {
		err := command.Signal(syscall.Signal(0))
		if err == nil {
			return fmt.Errorf("%s is running", *executableFile)
		}
	}
	return err
}

func Status() bool {
	err := GatherCommands()
	if err != nil {
		return false
	}
	if len(GatheredCommands) > 2 {
		return false
	}
	for _, command := range GatheredCommands {
		err := command.Signal(syscall.Signal(0))
		if err == nil {
			return true
		}
	}
	return false
}
