package telnet

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/vigilantsys/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "telnet",
	GadgetSynopsis: "telnet server and client functionality",
	GadgetUsage:    "telnet [-server] [-client] ADDRESS\n",
	Run:            Run,
	InitFlags:      initFlags,
}
var server bool
var client bool
var address string

func initFlags(f *flag.FlagSet) {
	flag.BoolVar(&server, "server", false, "enable server mode")
	flag.BoolVar(&client, "client", true, "enable client mode")
	flag.StringVar(&address, "add", "8080", "port to serve on, defaults to 8080")
}

func Run() {
	fmt.Printf("Starting server for directory %s on port %s\n\n", dir, port)
	fs := http.FileServer(http.Dir(dir))
	fmt.Println(http.ListenAndServe(":"+port, fs))
}
