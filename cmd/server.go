/*
Copyright © 2022 Vigilant Cyber Systems, Inc.
Sean Heath
<sheath@vigilantsys.com>
Marc Bohler
<mbohler@vigilantsys.com>
Dylan Harbaugh
<dharbaugh@vigilantsys.com>


Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "web server for sharing files",
	Long:  `Server launches a web server for sharing files. Files can be uploaded or downloaded from the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting server for directory %s on port %s\n", serverDirectory, serverPort)
		fmt.Println("/ for directory listing")
		fmt.Println("/upload for file upload")
		http.HandleFunc("/upload", uploadAction)
		fs := http.FileServer(http.Dir(serverDirectory))
		http.Handle("/", fs)
		fmt.Println(http.ListenAndServe(":"+serverPort, nil))
	},
}

var serverDirectory, serverPort string
var serverTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Upload File</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="http://{{ .host }}/upload"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>`

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverDirectory, "dir", ".", "the directory to serve files from")
	serverCmd.Flags().StringVar(&pivotPort, "port", "8080", "the port to listen on")

}

// Display the named template
func display(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	var networking = map[string]string{"host": r.Host}
	t, err := template.New("uploadTemplate").Parse(serverTemplate)
	if err != nil {
		panic(err)
	}
	t.Execute(w, networking)
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
	dst, err := os.Create(filepath.Clean(fmt.Sprintf("%s/%s", serverDirectory, handler.Filename)))
	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		return
	}
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
		display(w, r, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}
