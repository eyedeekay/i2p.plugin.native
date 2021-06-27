package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var GatheredCommands []*os.Process

func GatherCommands() error {
	err := filepath.Walk(*runDir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				pid, err := strconv.Atoi(info.Name())
				if err != nil {
					return nil
				}
				process, err := os.FindProcess(pid)
				if err != nil {
					Println("found stale pid in our management dir", pid)
					os.Remove(path)
				}
				Println("found old process in our management dir", pid)
				GatheredCommands = append(GatheredCommands, process)
			}
			return nil
		})
	return err
}

func StopCommands() error {
	err := GatherCommands()
	if err != nil {
		return err
	}
	for _, command := range GatheredCommands {
		Println("killing", command.Pid)
		if runtime.GOOS == "windows" {
			err = command.Signal(os.Kill)
		} else {
			err = command.Signal(os.Interrupt)
		}
		if err != nil {
			Println("warning:", err, command.Pid)
		}
		os.Remove(filepath.Join(*runDir, strconv.Itoa(command.Pid)))
	}
	return err
}
