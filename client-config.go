package shellservice

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ClientConfig struct {
	ClientName          *string
	ClientDisplayName   *string
	Command             *string
	CommandArgs         *string
	StopCommand         *string
	Delay               *string
	Start               *bool
	NoShellService      *bool
	CommandInPath       *bool
	Executable          *string
	ExtendClassPath     string
	JavaShellService    *string
	NoAutoSuffixWindows *bool
	TargetOS            *string
	ResourceDir         *string
}

func karenConfig() string {
	return ""
}

func (cc *ClientConfig) Print() string {
	r := "clientApp.0.main=" + *cc.JavaShellService + "\n"
	r += cc.PrintClientName()
	r += cc.PrintCommand()
	r += cc.PrintStop()
	r += cc.PrintDelay()
	r += cc.PrintStart()
	r += cc.PrintLibraries()
	r += karenConfig()
	return r
}

func (cc *ClientConfig) PrintLibraries() string {
	if cc.ExtendClassPath != "" {
		return fmt.Sprintf("clientApp.0.classpath=%s\n", cc.ExtendClassPath)
	}
	return ""
}

func (cc *ClientConfig) CheckClientName(name string) string {
	if cc.ClientName == nil || *cc.ClientName == "" {
		cc.ClientName = &name
	}
	return fmt.Sprintf("clientApp.0.name=%s\n", *cc.ClientName)
}

func (cc *ClientConfig) PrintClientName() string {
	if cc.ClientName == nil || *cc.ClientName == "" {
		log.Fatal("-name is a required field.")
	}
	return fmt.Sprintf("clientApp.0.name=%s\n", *cc.ClientName)
}

func (cc *ClientConfig) PrintCommandArgs() string {
	if cc.CommandArgs == nil || *cc.CommandArgs == "" {
		split := strings.Split(*cc.Command, " ")
		if len(split) > 1 {
			return strings.TrimRight(strings.Join(split[1:], " "), " ")
		}
	}
	return *cc.CommandArgs
}

func (cc *ClientConfig) PrintCommand() string {
	if cc.ClientName == nil || *cc.ClientName == "" {
		log.Fatal("-name is a required field.")
	}
	CIP := ""
	if cc.CommandInPath == nil || !*cc.CommandInPath {
		CIP = "$PLUGIN/lib/"
	}
	exesuffix := ""
	if *cc.TargetOS == "windows" && !*cc.NoAutoSuffixWindows {
		exesuffix = ".exe"
	}
	if cc.Command == nil || *cc.Command == "" {
		if strings.HasSuffix(*cc.Command, exesuffix) {
			exesuffix = ""
		}
		return fmt.Sprintf("clientApp.0.args=%s%s%s -shellservice.name \"%s\" -shellservice.displayname \"%s\" %s\n", CIP, *cc.Command, exesuffix, *cc.ClientName, *cc.ClientDisplayName, cc.PrintCommandArgs())
	}
	name := strings.Split(*cc.Command, " ")[0]
	if strings.HasSuffix(name, exesuffix) {
		exesuffix = ""
	}
	return fmt.Sprintf("clientApp.0.args=%s%s%s -shellservice.name \"%s\" -shellservice.displayname \"%s\" %s\n", CIP, name, exesuffix, *cc.ClientName, *cc.ClientDisplayName, cc.PrintCommandArgs())
}

func (cc *ClientConfig) PrintStop() string {
	return ""
}

func (cc *ClientConfig) PrintDelay() string {
	if cc.ClientName == nil || *cc.ClientName == "" {
		return fmt.Sprintf("clientApp.0.delay=%s\n", "5")
	}
	return fmt.Sprintf("clientApp.0.delay=%s\n", *cc.Delay)
}

func (cc *ClientConfig) PrintStart() string {
	if cc.Start == nil {
		return ""
	}
	return fmt.Sprintf("clientApp.0.startOnLoad=%t\n", *cc.Start)
}

func (cc *ClientConfig) CopyResDir() error {
	// TODO: move the copy resdir to a function
	// in client-config.go
	if cc.ResourceDir != nil && *cc.ResourceDir != "" {
		files := find(filepath.Join(*cc.ResourceDir, "lib"), ".jar")
		for i, file := range files {
			cleaned := strings.Replace(file, *cc.ResourceDir, "$PLUGIN/", 1)
			cc.ExtendClassPath += cleaned
			fmt.Printf("%d:%d-%s\n", i, len(files), cleaned)
			if i != len(files)-1 {
				cc.ExtendClassPath += ","
			}
		}
		if err := Copy(*cc.ResourceDir, "plugin/"); err != nil {
			return err
		}
	}
	return nil
}

func (cc *ClientConfig) CopyExecutable() error {
	exesuffix := ""
	if *cc.TargetOS == "windows" && !*cc.NoAutoSuffixWindows {
		if !strings.HasSuffix(*cc.Executable, ".exe") {
			exesuffix = ".exe"
		}
	}
	if err := Copy(*cc.Executable, "plugin/lib/"+*cc.Executable+exesuffix); err != nil {
		return err
	}
	if err := os.Chmod("plugin/lib/"+*cc.Executable+exesuffix, 0755); err != nil {
		return err
	}
	return nil
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
