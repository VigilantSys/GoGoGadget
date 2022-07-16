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

	"github.com/spf13/cobra"
	"github.com/reiver/go-telnet"
)


// telnetCmd represents the telnet command
var telnetCmd = &cobra.Command{
	Use:   "telnet",
	Short: "telnet client",
	Long:  `Telnet provides a telnet client program.`,
	Run: func(cmd *cobra.Command, args []string) {
		var caller telnet.Caller = telnet.StandardCaller

		err := telnet.DialToAndCall(fmt.Sprintf("%s:%s", telnetAddress, telnetPort), caller)
		if err != nil {
			fmt.Errorf("error connecting to Telnet server: %w", err)
		}

	},
}

var telnetAddress string
var telnetPort string

func init() {
	rootCmd.AddCommand(telnetCmd)
	telnetCmd.Flags().StringVar(&telnetAddress, "address", "", "address to connect to")
	telnetCmd.Flags().StringVar(&telnetPort, "port", "23", "port to connect to")
}

