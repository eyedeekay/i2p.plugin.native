package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	executableFile  = flag.String("exe", "", "Executable file to manage.")
	executableDir   = flag.String("exedir", "", "Directory to run the executable in.")
	executablePerms = flag.String("exeperm", "", "Change the permissions of the executable before running it.")
	executableArgs  = flag.String("args", "", "Pass a set of arguments to the executable")

	runDir          = flag.String("rundir", "", "Directory to store runtime files")
	restartAlways   = flag.Bool("restart-always", true, "If the managed application crashes, automatically restart it.")
	restartInterval = flag.String("restart-interval", "0m", "Periodically re-start the process Karen is managing")

	command    = flag.String("instruct", "", "Tell karen to 'start', 'stop', 'restart', or check the 'status' of a command running in the supervisor")
	commandLog = flag.String("log", "karen.log", "Redirect log to file")
	commandErr = flag.String("err", "karen.err", "Redirect error to file")

	verbose = flag.Bool("verbose", true, "Make karen talk a lot.")
)

var precheckErr error
var runErr error

func restartIntervalParsed() time.Duration {
	str := *restartInterval
	units := str[len(*restartInterval)-1:]
	num := strings.Replace(str, units, "", 1)
	delay, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	switch units {
	case "s":
		return time.Duration(delay) * time.Second
	case "m":
		return time.Duration(delay) * time.Minute
	case "h":
		return time.Duration(delay) * time.Hour
	case "d":
		return time.Duration(delay) * (time.Hour * time.Duration(24))
	default:
		return time.Duration(delay) * time.Minute
	}
}

func checkRestartInterval() {
	rip := restartIntervalParsed()
	if rip.Seconds() != 0 {
		for {
			time.Sleep(rip)
			runErr = Stop()
			checkRunErr()
			if !(Status()) {
				runErr = Run()
				checkRunErr()
			}
		}
	}
}

func checkExecutable() {
	if *executableFile == "" {
		log.Fatal("-exe is a required field")
	}
}

func checkExecutableDir() {
	if *executableDir == "" {
		preDir, precheckErr := filepath.Abs(*executableFile)
		if precheckErr != nil {
			log.Fatal("Run directory not provided and absolute path cannot be determined.")
		}
		*executableDir = filepath.Dir(preDir)
	}
}

func getExecutable() string {
	return filepath.Join(*executableDir, filepath.Base(*executableFile))
}

func getArgs() []string {
	args := strings.Split(*executableArgs, " ")
	return args
}

func checkExecutablePerms() {
	Println("Checking permissions")
	if runtime.GOOS != "windows" {
		if *executablePerms == "" {
			return
		}
		file, precheckErr := os.Open(getExecutable())
		if precheckErr != nil {
			log.Fatal(precheckErr)
		}
		defer file.Close()
		tempval, precheckErr := strconv.ParseUint(*executablePerms, 8, 32)
		if precheckErr != nil {
			log.Fatal(precheckErr)
		}
		Printf("fixing permissions on %s to %o", getExecutable(), tempval)
		err := file.Chmod(os.FileMode(tempval))
		if err != nil {
			ioutil.WriteFile(filepath.Join(*executableDir, "perm.err"), []byte(err.Error()), 0644)
		}
	}
}

func checkRunDir() {
	if *runDir == "" {
		*runDir = filepath.Join(filepath.Dir(*executableDir), "run")
	}
	os.MkdirAll(*runDir, 0755)
}

func checkRunErr() {
	if runErr != nil {
		log.Fatal(runErr)
	}
}

func main() {
	flag.Parse()
	checkExecutable()
	checkExecutableDir()
	checkExecutablePerms()
	checkRunDir()
	Printf(
		"%s -exe %s \n\t-exedir %s \n\t-exeperm %v \n\t-args %s \n\t-rundir %s \n\t-restart-always %t \n\t-instruct %s \n\t-verbose %t",
		os.Args[0],
		*executableFile,
		*executableDir,
		*executablePerms,
		*executableArgs,
		*runDir,
		*restartAlways,
		*command,
		*verbose,
	)
	runErr = GatherCommands()
	checkRunErr()

	if _, err := os.Stat(filepath.Join(*executableDir, "karen.pid")); os.IsNotExist(err) {
		ioutil.WriteFile(filepath.Join(*executableDir, "karen.pid"), []byte(strconv.Itoa(os.Getpid())), 0644)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		err := Stop()
		removePid()
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()
	go checkRestartInterval()

	defer Stop()
	defer removePid()
	for {
		switch *command {
		case "start":
			if !(Status()) {
				runErr = Run()
				checkRunErr()
				WritePID()
			}
		case "stop":
			runErr = Stop()
			checkRunErr()
			*restartAlways = false
		case "restart":
			runErr = Stop()
			checkRunErr()
			if !(Status()) {
				runErr = Run()
				checkRunErr()
			}
		case "status":
			log.Println(StatusMessage())
			os.Exit(0)
		default:
			log.Println("Please choose to start, stop, restart, or check the status of your command")
			os.Exit(1)
		}
		if cmd != nil {
			cmd.Wait()
		}
		if !*restartAlways {
			os.Exit(0)
		}
	}
}

func removePid() {
	os.Remove(filepath.Join(*executableDir, "karen.pid"))
}
