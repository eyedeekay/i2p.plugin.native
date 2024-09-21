package shellservice

import (
	"fmt"
	//"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	ClientName          *string `yaml:"clientName,omitempty"`
	ClientDisplayName   *string `yaml:"clientDisplayName,omitempty"`
	Command             *string `yaml:"command,omitempty,omitempty"`
	CommandArgs         *string `yaml:"commandArgs,omitempty"`
	StopCommand         *string `yaml:"stopCommand,omitempty"`
	Delay               *string `yaml:"delay,omitempty"`
	Start               *bool   `yaml:"start,omitempty"`
	NoShellService      *bool   `yaml:"noShellService,omitempty"`
	CommandInPath       *bool   `yaml:"commandInPath,omitempty"`
	Executable          *string `yaml:"executable,omitempty"`
	ExtendClassPath     string  `yaml:"-,omitempty"`
	JavaShellService    *string `yaml:"javaShellService,omitempty"`
	NoAutoSuffixWindows *bool   `yaml:"noAutoSuffixWindows,omitempty"`
	TargetOS            *string `yaml:"targetOS,omitempty"`
	ResourceDir         *string `yaml:"resourceDir,omitempty"`
	I2PTunnelConfig     *string `yaml:"i2ptunnelConfig,omitempty"`
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
	r += cc.PrintI2PTunnelClientApp()
	r += cc.PrintI2PTunnelClientName()
	r += cc.PrintI2PTunnelClientArgs()
	r += cc.PrintI2PTunnelClientDelay()
	r += cc.PrintI2PTunnelClientStart()
	r += karenConfig()
	return Replace(r)
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

func (cc *ClientConfig) PrintI2PTunnelClientApp() string {
	if cc.I2PTunnelConfig == nil || *cc.I2PTunnelConfig == "" {
		return ""
	}
	return fmt.Sprintf("clientApp.2.main=%s\n", "net.i2p.i2ptunnel.TunnelControllerGroup")
}

func (cc *ClientConfig) PrintI2PTunnelClientName() string {
	if cc.I2PTunnelConfig == nil || *cc.I2PTunnelConfig == "" {
		return ""
	}
	return fmt.Sprintf("clientApp.1.name=%s\n", *cc.ClientName+"i2ptunnel")
}

func (cc *ClientConfig) PrintI2PTunnelClientArgs() string {
	if cc.I2PTunnelConfig == nil || *cc.I2PTunnelConfig == "" {
		return ""
	}
	return fmt.Sprintf("clientApp.1.args=$PLUGIN/%s\n", *cc.I2PTunnelConfig)
}

func (cc *ClientConfig) PrintI2PTunnelClientDelay() string {
	if cc.I2PTunnelConfig == nil || *cc.I2PTunnelConfig == "" {
		return ""
	}
	return fmt.Sprintf("clientApp.1.delay=%s\n", "-1")
}

func (cc *ClientConfig) PrintI2PTunnelClientStart() string {
	if cc.I2PTunnelConfig == nil || *cc.I2PTunnelConfig == "" {
		return ""
	}
	return fmt.Sprintf("clientApp.1.startOnLoad=%s\n", "true")
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
			return fmt.Errorf("CopyResDir(): %s", err)
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
		return fmt.Errorf("CopyExecutable(): %s", err)
	}
	if err := os.Chmod("plugin/lib/"+*cc.Executable+exesuffix, 0755); err != nil {
		return fmt.Errorf("CopyExecutable(): %s", err)
	}
	return nil
}

func (cc *ClientConfig) Load() error {
	if _, err := os.Stat(clientFile()); os.IsNotExist(err) {
		return nil
	}
	yamlFile, err := ioutil.ReadFile(clientFile())
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, cc)
}

func (cc *ClientConfig) Save() error {
	bytes, err := yaml.Marshal(cc)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(clientFile(), bytes, 0644)
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
	return copy.Copy(src, dst)
}

func clientFile() string {
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")
	r := "client"
	if goos != "" {
		r += "-" + goos
	}
	if goarch != "" {
		r += "-" + goarch
	}
	return r + ".yaml"
}
