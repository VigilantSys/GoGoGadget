package server

import (
	"flag"
	"fmt"
	"os"
	"io"
	"net/http"
	"text/template"
	"path/filepath"

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

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("./upload.html"))

func initFlags(f *flag.FlagSet) {
	f.StringVar(&dir, "dir", ".", "directory to serve, defaults to current directory")
	f.StringVar(&port, "port", "8080", "port to serve on, defaults to 8080")
}

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(filepath.Clean(fmt.Sprintf("%s/%s", dir, handler.Filename)))
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func uploadAction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			display(w, "upload", nil)
		case "POST":
			uploadFile(w, r)
	}
}

func Run() {
	fmt.Printf("Starting server for directory %s on port %s\n\n", dir, port)
	http.HandleFunc("/upload", uploadAction)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)
	fmt.Println(http.ListenAndServe(":"+port, nil))
}
