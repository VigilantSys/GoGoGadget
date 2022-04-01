package download

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/vigilantsys/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "download",
	GadgetSynopsis: "like curl - download a file",
	GadgetUsage:    "download\n",
	Run:            Run,
	InitFlags:      initFlags,
}
var url string
var outfile string
var insecure bool
var rtype string

func initFlags(f *flag.FlagSet) {
	f.StringVar(&url, "url", "", "URL to download")
	f.StringVar(&outfile, "outfile", "outfile", "File to write to")
	f.BoolVar(&insecure, "k", false, "Skip certificate verification")
	f.StringVar(&rtype, "request", "GET", "Select request type (GET, POST, etc...)")
}

func Run() {
	if url == "" {
		fmt.Println("download error - no URL provided")
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(rtype, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: error checking
	// TODO: check for blank url
	// TODO: status bar
	resp, err := client.Do(req)
	if err != nil {
		// TODO: do better.
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Make the file
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
