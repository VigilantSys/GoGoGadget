//go:build windows || android || (linux && (arm || amd64 || arm64 || 386))

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
	"os"
	"time"
	"image/png"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kbinani/screenshot"

)

var outputDir = ""
// screenshotCmd represents the screenshot command
var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		n := screenshot.NumActiveDisplays()

		for i :=0; i < n; i++ {
			bounds := screenshot.GetDisplayBounds(i)
			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				fmt.Println("error taking screenshot: " + err.Error())
				return
			}
			ts := time.Now().Format("20060102-15:04:05")
			if outputDir != "" {
				// Check for trailing slash
				if !strings.HasSuffix(outputDir, "/") {
					// no slash found, add one
					outputDir = outputDir + "/"
				}
				// prepend directory to timestamp
				ts = outputDir + ts 
			}
			fname := fmt.Sprintf("%s-%d-%dx%d.png", ts, i, bounds.Dx(), bounds.Dy())
			file, err := os.Create(fname)
			if err != nil {
				fmt.Println("could not save files to the specified location: " + fname)
				fmt.Println("error: " + err.Error())
				return
			}
			defer file.Close()
			png.Encode(file, img)
			fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fname)
		}
	},
}

func init() {
	rootCmd.AddCommand(screenshotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// screenshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// screenshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	screenshotCmd.Flags().StringVarP(&outputDir, "outdir", "o", "", "directory to write images to")
}
