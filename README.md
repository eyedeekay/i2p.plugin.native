ShellService Plugin Generator
=============================

Generates a valid ShellService plugin from a single script or executable. A ShellService plugin is an application which runs as a process which is managed by the I2P sofware. This allows I2P to manage the lifetime of a non-JVM application by monitoring it. That way, plugins don't risk outliving the I2P router because the router loses track of them.

A ShellService plugin should not "daemonize" itself or othewise fork itself to the background. The I2P router will manage it as a "Client Application" and daemonizing it will break this. If your process forks itself to the background by default, it must have this feature disabled to work as a ShellService.

ShellService adds additional "arguments" to the applications in question, which are used to configure the ShellService. In order to function correctly, the ShellService must be named so that the name of the Plugin, i.e. name in plugin.config, must match the -shellservice.name argument added by the ShellService class. If it does not, the ShellService will be unable to look up the plugin process. If the final elements of the ShellService name are -$OS-$ARCH or -$OS they may be optionally omitted.
There are some examples in the Makefile for now.

Here's a copy of the usage while I work on a better README.md:

```markdown
Usage of i2p.plugin.native:
  -author string
    	Author
  -autostart
    	Start client automatically (default true)
  -clientname string
    	Name of the client, defaults to same as plugin
  -command string
    	Command to start client, defaults to $PLUGIN/lib/exename
  -commandargs string
    	Pass arguments to command
  -consoleicon string
    	Icon to use in console for Web Apps only. Use icondata for native apps.
  -consolename string
    	Name to use in the router console sidebar
  -consoleurl string
    	URL to use in the router console sidebar
  -date string
    	Release Date
  -delaystart string
    	Delay start of client by seconds (default "1")
  -desc string
    	Plugin description
  -exename string
    	Name of the executable the plugin will run, defaults to name
  -icondata string
    	Path to icon for console, which i2p.plugin.native will automatically encode
  -installonly
    	Only allow installing with this plugin, fail if a previous installation exists
  -javashellservice string
    	specify ShellService java path (default "net.i2p.router.web.ShellService")
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
  -noautosuffixwindows
    	Don't automatically add .exe to exename on Windows
  -nostart
    	Don't automatically start the plugin after installing
  -nostop
    	Disable stopping the plugin from the console
  -pathcommand
    	Wrap a command found in the system $PATH, don't prefix the command with $PLUGIN/lib/
  -res string
    	a directory of additional resources to include in the plugin
  -restart
    	Require a router restart after installing or updating the plugin
  -signer string
    	Signer of the plugin
  -signer-dir string
    	Directory to look for signing keys
  -stopcommand string
    	Command to stop client, defaults to killall exename
  -targetos string
    	Target to run the plugin on
  -updateonly
    	Only allow updates with this plugin, fail if no previous installation exists
  -updateurl string
    	The URL to retrieve updates from, defaults to website+pluginname.su3
  -version string
    	Version of the plugin
  -website string
    	The website of the plugin
```
