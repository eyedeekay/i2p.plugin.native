package main

import (
	"fmt"
	"log"
	"strings"
)

type ClientConfig struct {
	ClientName  *string
	Command     *string
	CommandArgs *string
	StopCommand *string
	Delay       *string
	Start       *bool
}

func karenConfig() string {
	return ""
}

func (cc *ClientConfig) Print() string {
	r := "clientApp.0.main=net.i2p.util.ShellService\n"
	r += cc.PrintClientName()
	r += cc.PrintCommand()
	r += cc.PrintStop()
	r += cc.PrintDelay()
	r += cc.PrintStart()
	r += karenConfig()
	return r
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
	if cc.Command == nil || *cc.Command == "" {
		if cc.ClientName != nil || *cc.ClientName != "" {
			return fmt.Sprintf("clientApp.0.args=$PLUGIN/lib/%s -shellservice.name=\"%s\" -shellservice.displayname=\"%s\" %s\n", *cc.Command, *cc.ClientName, *cc.ClientName, cc.PrintCommandArgs())
		}
		log.Fatal("-name is a required field.")
	}
	name := strings.Split(*cc.Command, " ")[0]
	return fmt.Sprintf("clientApp.0.args=$PLUGIN/lib/%s -shellservice.name=\"%s\" -shellservice.displayname=\"%s\" %s\n", name, name, name, cc.PrintCommandArgs())
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