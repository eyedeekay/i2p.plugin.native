package main

import (
	"fmt"
	"log"
)

type ClientConfig struct {
	ClientName  *string
	Command     *string
	StopCommand *string
	Delay       *string
	Start       *bool
}

func (cc *ClientConfig) Print() string {
	r := "clientApp.0.main=net.i2p.util.ShellCommand\n"
	r += cc.PrintClientName()
	r += cc.PrintCommand()
	r += cc.PrintStop()
	r += cc.PrintDelay()
	r += cc.PrintStart()
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

func (cc *ClientConfig) PrintCommand() string {
	if cc.Command == nil || *cc.Command == "" {
		if cc.ClientName != nil || *cc.ClientName != "" {
			return fmt.Sprintf("clientApp.0.args=$PLUGIN/lib/%s 2>&1 | tee $PLUGIN/railroad.log\n", *cc.ClientName)
		}
		log.Fatal("-name is a required field.")
	}
	return fmt.Sprintf("clientApp.0.args=%s\n", *cc.Command)
}

func (cc *ClientConfig) PrintStop() string {
	if cc.StopCommand == nil || *cc.StopCommand == "" {
		return fmt.Sprintf("clientApp.0.stopargs=killall %s\n", *cc.ClientName)
	}
	return fmt.Sprintf("clientApp.0.stopargs=%s\n", *cc.StopCommand)
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
