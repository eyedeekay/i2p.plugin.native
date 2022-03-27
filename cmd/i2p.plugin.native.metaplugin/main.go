package main

import (
	//	"archive/zip"

	//	"runtime"

	"archive/zip"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "i2pgit.org/idk/i2p.plugin.native"
	"i2pgit.org/idk/reseed-tools/su3"
)

var pc PluginConfig
var cc ClientConfig

// The goal here is to create a plugin which can be installed which has no contents, but contains update
// information for other plugins, i.e.
//
// i2p.plugins.tor-manager.su3
//
// would contain meta-information that matches the $OS and $ARCH of the runtime machine
// which would then update at the first opportunity, downloading the latest version of the plugin
// for the appropriate architecture with the correct binaries. It doesn't save me any work :) but
// it should be nice for the user to be able to install the plugin without having to comb through
// OS/ARCH pairs they don't really give a damn about.

var wd, wderr = os.Getwd()
var dir = flag.String("dir", wd, "Directory to look for plugins in")

func flagSet() {
	if wderr != nil {
		log.Fatal(wderr)
	}
	flag.Parse()
}

func FindSu3sInDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var su3s []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".su3") {
			su3s = append(su3s, f.Name())
		}
	}
	return su3s, nil
}

func MatchSu3sWithOSARCH(su3s []string) []string {
	var matched []string
	for _, su3 := range su3s {
		for _, osarch := range ARCHES {
			if strings.Contains(su3, osarch) {
				matched = append(matched, su3)
			}
		}
		for _, osarch := range OSES {
			if strings.Contains(su3, osarch) {
				matched = append(matched, su3)
			}
		}
	}
	return removeDuplicateValues(matched)
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	flagSet()
	absdir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}
	su3s, err := FindSu3sInDir(absdir)
	if err != nil {
		log.Fatal(err)
	}
	su3s = MatchSu3sWithOSARCH(su3s)
	// extract the plugin.config and the clients.config from topmost su3
	// and get the plugin name from the plugin.config file.
	// copy the plugin.config and clients.config to meta/plugin.config and meta/clients.config
	// create a script that says "this is a meta plugin, downloading an update for the plugin
	// will install a real version" in each location per-platform, i.e. plugin-name-os-arch
	// use sh for Linux and OSX and bat for Windows
	su3file := su3.New()
	bytes, err := ioutil.ReadFile(filepath.Join(absdir, su3s[0]))
	if err != nil {
		log.Fatal(err)
	}
	err = su3file.UnmarshalBinary(bytes)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(su3s[0]+".zip", su3file.Content, 0644)
	zipReader, err := zip.OpenReader(su3s[0] + ".zip")
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()
	clientConfig, err := zipReader.Open("clients.config")
	if err != nil {
		log.Fatal(err)
	}
	defer clientConfig.Close()
	pluginConfig, err := zipReader.Open("plugin.config")
	if err != nil {
		log.Fatal(err)
	}
	defer pluginConfig.Close()
	clientConfigBytes, err := ioutil.ReadAll(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	pluginConfigBytes, err := ioutil.ReadAll(pluginConfig)
	if err != nil {
		log.Fatal(err)
	}
	os.Mkdir("meta", 0755)
	replacedClientConfigBytes := []byte(Replace(string(clientConfigBytes)))
	replacedPluginConfigBytes := []byte(Replace(string(pluginConfigBytes)))
	err = ioutil.WriteFile("meta/clients.config", replacedClientConfigBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("meta/plugin.config", replacedPluginConfigBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
