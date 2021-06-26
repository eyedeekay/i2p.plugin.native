I2P native plugin generation tool
=================================

I wrote this way faster than I documented it. Shocking, right?

This is a handy little tool for assembling I2P plugins when those
plugins don't have a clean way to interface with the JVM, or just don't
need one. Think of it a little like `checkinstall` but for I2P Plugins.
Right now it mostly works, and it's pretty cleanly put together.

There are some examples in the Makefile for now.

Here's a copy of the usage while I work on a better README.md:

```bash
Usage of ./scripts/bin/i2p.plugin.native:
  -author string
    	Author
  -autostart
    	Start client automatically (default true)
  -clientname string
    	Name of the client, defaults to same as plugin
  -command string
    	Command to start client, defaults to $PLUGIN/exename
  -consoleicon string
    	Icon to use in console for Web Apps only. Use icondata for native apps.
  -consolename string
    	Name to use in the router console sidebar
  -date string
    	Release Date
  -delaystart string
    	Delay start of client by seconds (default "5")
  -desc string
    	Plugin description
  -exename string
    	Name of the executable the plugin will run, defaults to name
  -icondata string
    	Path to icon for console, which i2p.plugin.native will automatically encode
  -installonly
    	Only allow installing with this plugin, fail if a previous installation exists
  -key string
    	Key to use(omit for su3)
  -license string
    	License of the plugin
  -max string
    	Maximum I2P version
  -max-jetty string
    	Maximum Jetty version
  -min string
    	Minimum I2P version
  -min-java string
    	Minimum Java version
  -min-jetty string
    	Minimum Jetty version
  -name string
    	Name of the plugin
  -nostart
    	Don't automatically start the plugin after installing
  -nostop
    	Disable stopping the plugin from the console
  -res string
    	a directory of additional resources to include in the plugin
  -restart
    	Require a router restart after installing or updating the plugin
  -signer string
    	Signer of the plugin
  -stopcommand string
    	Command to stop client, defaults to killall exename
  -updateonly
    	Only allow updates with this plugin, fail if no previous installation exists
  -updateurl string
    	The URL to retrieve updates from, defaults to website+pluginname.su3
  -version string
    	Version of the plugin
  -website string
    	The website of the plugin
```
