package ggg-telnet

import (
	"flag"

	"github.com/reiver/go-telnet"
	"github.com/vigilantsys/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "telnet",
	GadgetSynopsis: "telnet server and client functionality",
	GadgetUsage:    "telnet [-server] [-client] ADDRESS\n",
	Run:            Run,
	InitFlags:      initFlags,
}
var serverPort int
var clientAddress string

func initFlags(f *flag.FlagSet) {
	flag.IntVar(&serverPort, "server", 23, "enable server mode")
	flag.StringVar(&clientAddress, "client", "", "port to serve on, defaults to 8080")
}

func Run() {
	if clientAddress != "" {
		var caller telnet.Caller = telnet.StandardCaller
		telent.DialToAndCall(clientAddress, caller)
	}
}
