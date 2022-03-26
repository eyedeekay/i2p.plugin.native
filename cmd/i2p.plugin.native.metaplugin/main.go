package main

import (
	//	"archive/zip"

	//	"runtime"

	. "i2pgit.org/idk/i2p.plugin.native"
)

var pc PluginConfig
var cc ClientConfig

// The goal here is to create a plugin which can be installed which has no contents, but contains update
// information for other plugins, i.e.
//
// i2p.plugins.tor-manager-amd64
//
// would contain meta-information that matches the $OS and $ARCH of the runtime machine
// which would then update at the first opportunity, downloading the latest version of the plugin
// for the appropriate architecture with the correct binaries. It doesn't save me any work :) but
// it should be nice for the user to be able to install the plugin without having to comb through
// OS/ARCH pairs they don't really give a damn about.

func main() {

}
