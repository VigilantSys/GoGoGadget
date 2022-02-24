package server

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/seandheath/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "server",
	GadgetSynopsis: "web server to share files in a directory",
	GadgetUsage:    "server\n",
	Run:            Run,
	InitFlags:      initFlags,
}
var dir string
var port string

func initFlags(f *flag.FlagSet) {
	flag.StringVar(&dir, "dir", ".", "directory to serve, defaults to current directory")
	flag.StringVar(&port, "port", "8080", "port to serve on, defaults to 8080")
}

func Run() {
	fmt.Printf("Starting server for directory %s on port %s\n\n", dir, port)
	fs := http.FileServer(http.Dir(dir))
	fmt.Println(http.ListenAndServe(":"+port, fs))
}
