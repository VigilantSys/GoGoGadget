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
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "like curl - download a file",
	Long: `Download a file via HTTP request and save it to the target destination.
	
Example:
  gogogadget download -url http://www.google.com -outfile google.html`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("download error - no URL provided")
			os.Exit(1)
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}

		client := &http.Client{Transport: tr}
		req, err := http.NewRequest(request, url, nil)
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
	},
}

var url, outfile, request string
var insecure bool

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	downloadCmd.Flags().StringVar(&url, "url", "", "URL to download")
	downloadCmd.Flags().StringVar(&outfile, "outfile", "outfile", "File to write to")
	downloadCmd.Flags().StringVar(&request, "request", "GET", "Select request type (GET, POST, etc...)")
	downloadCmd.Flags().BoolVar(&insecure, "k", false, "Skip certificate verification")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
