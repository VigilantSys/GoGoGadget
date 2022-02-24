package pivot

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/seandheath/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "pivot",
	GadgetSynopsis: "pivot traffic in one port and out to a target",
	GadgetUsage:    "pivot\n",
	Run:            Run,
	InitFlags:      initFlags,
}
var target string
var port string
var protocol string

func initFlags(f *flag.FlagSet) {
	f.StringVar(&target, "target", "", "host and port to forward traffic to <host>:<port>")
	f.StringVar(&port, "port", "8080", "local port to listen on")
	f.StringVar(&protocol, "protocol", "tcp", "protocol to use, defaults to tcp")
}

func Run() {
	// TODO: Check for required variables
	// Start Server
	incoming, err := net.Listen(protocol, fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done with listen")
	fmt.Printf("server running %s\n", port)

	// Accept Connection
	for {
		client, err := incoming.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(client, target, protocol)
	}
}

func handleRequest(client net.Conn, targetAddress string, protocol string) {
	fmt.Printf("client '%v' connected!\n", client.RemoteAddr())

	// Dial out to the target

	target, err := net.Dial(protocol, targetAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connection established to %v\n", target.RemoteAddr())
	go CopyIO(client, target)
	go CopyIO(target, client)
}

func CopyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
