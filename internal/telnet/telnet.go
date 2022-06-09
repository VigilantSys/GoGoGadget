package telnet

import (
	"flag"

	"github.com/vigilantsys/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "telnet",
	GadgetSynopsis: "telnet client functionality",
	GadgetUsage:    "telnet ADDRESS\n",
	Run:            Run,
	InitFlags:      initFlags,
}

var addr string

func initFlags(f *flag.FlagSet) {
	flag.StringVar(&addr, "address", "", "address for telnet server to connect to")
}

func Run() {
}
