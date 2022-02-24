package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/seandheath/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget
var dir string
var port string

func init() {
	Gadget = gadget.New("server", Run)

	if len(os.Args) < 2 || os.Args[1] != "server" {
		return
	}
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&dir, "dir", path, "directory to serve, defaults to current directory")
	flag.StringVar(&port, "port", "8080", "port to serve on, defaults to 8080")
}

func Run(args []string) {
	fmt.Printf("Starting server for directory %s on port %s\n\n", dir, port)
	fs := http.FileServer(http.Dir(dir))
	fmt.Println(http.ListenAndServe(":"+port, fs))
}
