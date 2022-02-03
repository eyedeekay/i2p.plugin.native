package main

import (
	//	"archive/zip"

	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	//	"runtime"

	. "i2pgit.org/idk/i2p.plugin.native"

	"i2pgit.org/idk/reseed-tools/su3"
)

var pc PluginConfig
var cc ClientConfig

//var executable string
//var resdir *string

//var targetos *string
//var noautosuffixwindows *bool

var javaShellService = "net.i2p.router.web.ShellService"

func flagsSet() {
	pc.PluginName = flag.String("name", "", "Name of the plugin")
	pc.KeyName = flag.String("key", "", "Key to use(omit for su3)")
	pc.Signer = flag.String("signer", "", "Signer of the plugin")
	pc.SignerDirectory = flag.String("signer-dir", "", "Directory to look for signing keys")
	pc.Version = flag.String("version", "", "Version of the plugin")
	pc.License = flag.String("license", "", "License of the plugin")
	pc.Date = flag.String("date", "", "Release Date")
	pc.Author = flag.String("author", "", "Author")
	pc.Website = flag.String("website", "", "The website of the plugin")
	pc.UpdateURL = flag.String("updateurl", "", "The URL to retrieve updates from, defaults to website+pluginname.su3")
	pc.Description = flag.String("desc", "", "Plugin description")
	//pc.DescriptionLang     []flag.String("","","")
	pc.ConsoleLinkName = flag.String("consolename", "", "Name to use in the router console sidebar")
	//pc.ConsoleLinkNameLang []flag.String("","","")
	pc.ConsoleIcon = flag.String("consoleicon", "", "Icon to use in console for Web Apps only. Use icondata for native apps.")
	pc.ConsoleIconCode = flag.String("icondata", "", "Path to icon for console, which i2p.plugin.native will automatically encode")
	pc.ConsoleLinkURL = flag.String("consoleurl", "", "URL to use in the router console sidebar")
	pc.MinVersion = flag.String("min", "", "Minimum I2P version")
	pc.MaxVerion = flag.String("max", "", "Maximum I2P version")
	pc.MinJava = flag.String("min-java", "", "Minimum Java version")
	pc.MinJetty = flag.String("min-jetty", "", "Minimum Jetty version")
	pc.MaxJetty = flag.String("max-jetty", "", "Maximum Jetty version")
	pc.NoStop = flag.Bool("nostop", false, "Disable stopping the plugin from the console")
	pc.NoStart = flag.Bool("nostart", false, "Don't automatically start the plugin after installing")
	pc.Restart = flag.Bool("restart", false, "Require a router restart after installing or updating the plugin")
	pc.OnlyUpdate = flag.Bool("updateonly", false, "Only allow updates with this plugin, fail if no previous installation exists")
	pc.OnlyInstall = flag.Bool("installonly", false, "Only allow installing with this plugin, fail if a previous installation exists")
	cc.ClientName = flag.String("clientname", "", "Name of the client, defaults to same as plugin")
	cc.Command = flag.String("command", "", "Command to start client, defaults to $PLUGIN/lib/exename")
	cc.CommandArgs = flag.String("commandargs", "", "Pass arguments to command")
	cc.StopCommand = flag.String("stopcommand", "", "Command to stop client, defaults to killall exename")
	cc.Delay = flag.String("delaystart", "1", "Delay start of client by seconds")
	cc.Start = flag.Bool("autostart", true, "Start client automatically")
	//cc.NoShellService = flag.Bool("noshellservice", false, "Use ShellCommand+Karen instead of ShellService to generate plugin")
	cc.CommandInPath = flag.Bool("pathcommand", false, "Wrap a command found in the system $PATH, don't prefix the command with $PLUGIN/lib/")
	cc.Executable = flag.String("exename", "", "Name of the executable the plugin will run, defaults to name")
	cc.ResourceDir = flag.String("res", "", "a directory of additional resources to include in the plugin")
	cc.TargetOS = flag.String("targetos", os.Getenv("GOOS"), "Target to run the plugin on")
	cc.NoAutoSuffixWindows = flag.Bool("noautosuffixwindows", false, "Don't automatically add .exe to exename on Windows")
	cc.JavaShellService = flag.String("javashellservice", javaShellService, "specify ShellService java path")
	flag.Parse()
	cc.ClientDisplayName = pc.ConsoleLinkName
}

func goBin() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return filepath.Join(gopath, "bin")
}

func main() {
	flagsSet()

	if *cc.Executable != "" {
		*cc.ClientName = *cc.Executable
	}

	cc.CheckClientName(*pc.PluginName)

	if *cc.Executable == "" {
		*cc.Executable = *cc.ClientName
	}

	fmt.Printf("executable:%s\n", *cc.Executable)
	fmt.Printf("resources:%s\n", *cc.ResourceDir)

	os.RemoveAll("plugin")
	if err := os.MkdirAll("plugin/lib", 0755); err != nil {
		log.Fatal(err)
	}
	if err := cc.CopyResDir(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("plugin.config:\n\t%s\n", pc.Print())
	fmt.Printf("client.config:\n\t%s\n", cc.Print())

	if err := ioutil.WriteFile("plugin/plugin.config", []byte(pc.Print()), 0644); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("plugin/clients.config", []byte(cc.Print()), 0644); err != nil {
		log.Fatal(err)
	}

	// TODO: move exe-copy logic into plugin-config.go
	if err := cc.CopyExecutable(); err != nil {
		log.Fatal(err)
	}

	if err := createZip(); err != nil {
		log.Fatal(err)
	}
	if file, err := createSu3(); err != nil {
		log.Fatal(err)
	} else {
		if data, err := file.MarshalBinary(); err != nil {
			log.Fatal(err)
		} else {
			ioutil.WriteFile(*pc.PluginName+".su3", data, 0644)
		}
	}
}

func createZip() error {
	return pc.CreateZip()
}

func createSu3() (*su3.File, error) {
	return pc.CreateSu3()
}
