package main

import (
	"archive/zip"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"

	"i2pgit.org/idk/reseed-tools/su3"
)

var pc PluginConfig
var cc ClientConfig

var executable string
var resdir *string

func flagsSet() {
	pc.PluginName = flag.String("name", "", "Name of the plugin")
	pc.KeyName = flag.String("key", "", "Key to use(omit for su3)")
	pc.Signer = flag.String("signer", "", "Signer of the plugin")
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
	cc.Command = flag.String("command", "", "Command to start client, defaults to $PLUGIN/exename")
	cc.StopCommand = flag.String("stopcommand", "", "Command to stop client, defaults to killall exename")
	cc.Delay = flag.String("delaystart", "5", "Delay start of client by seconds")
	cc.Start = flag.Bool("autostart", true, "Start client automatically")
	executable = *flag.String("exename", "", "Name of the executable the plugin will run, defaults to name")
	resdir = flag.String("res", "config", "a directory of additional resources to include in the plugin")
	flag.Parse()
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

func main() {
	flagsSet()

	if executable != "" {
		cc.ClientName = &executable
	}

	cc.CheckClientName(*pc.PluginName)

	if executable == "" {
		executable = *cc.ClientName
	}

	fmt.Printf("executable:%s\n", executable)
	fmt.Printf("resources:%s\n", resdir)
	fmt.Printf("plugin.config:\n\t%s\n", pc.Print())
	fmt.Printf("client.config:\n\t%s\n", cc.Print())

	if err := os.MkdirAll("plugin/lib", 0755); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("plugin/plugin.config", []byte(pc.Print()), 0644); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("plugin/clients.config", []byte(cc.Print()), 0644); err != nil {
		log.Fatal(err)
	}

	if err := Copy(executable, "plugin/lib/"+executable); err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod("plugin/lib/"+executable, 0755); err != nil {
		log.Fatal(err)
	}
	if resdir != nil && *resdir != "" {
		if err := copy.Copy(*resdir, "plugin/lib/"); err != nil {
			log.Fatal(err)
		}
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
	file, err := os.Create(*pc.PluginName + ".zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		zippath := strings.Replace(path, "plugin/", "", 1)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(zippath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk("plugin", walker)
	return err
}

func createSu3() (*su3.File, error) {
	su3File := su3.New()
	su3File.FileType = su3.FileTypeZIP
	su3File.ContentType = su3.ContentTypePlugin

	err := createZip()
	if err != err {
		return nil, err
	}
	zipped, err := ioutil.ReadFile(*pc.PluginName + ".zip")
	if err != err {
		return nil, err
	}
	su3File.Content = zipped

	su3File.SignerID = []byte(*pc.Signer)
	sk, err := loadPrivateKey(*pc.Signer)
	if err != nil {
		return nil, err
	}
	su3File.Sign(sk)
	return su3File, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privPem, err := ioutil.ReadFile(strings.Replace(path, "@", "_at_", -1) + ".pem")
	if err != nil {
		return nil, err
	}

	privDer, _ := pem.Decode(privPem)
	privKey, err := x509.ParsePKCS1PrivateKey(privDer.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
