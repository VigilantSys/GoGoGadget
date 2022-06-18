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
	"log"
	"net"

	"github.com/spf13/cobra"
)

// pivotCmd represents the pivot command
var pivotCmd = &cobra.Command{
	Use:   "pivot",
	Short: "pivot traffic in one port and out to a target",
	Long:  `Pivot accepts traffic in one port and forwards it to the provided destination.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Check for required variables
		// Start Server
		incoming, err := net.Listen(pivotProtocol, fmt.Sprintf(":%s", pivotPort))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("done with listen")
		fmt.Printf("server running %s\n", pivotPort)

		// Accept Connection
		for {
			client, err := incoming.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go handleRequest(client, pivotTarget, pivotProtocol)
		}
	},
}

var pivotTarget, pivotPort, pivotProtocol string

func init() {
	rootCmd.AddCommand(pivotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pivotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pivotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pivotCmd.Flags().StringVar(&pivotTarget, "target", "", "the target host to forward traffic to <host>:<port>")
	pivotCmd.Flags().StringVar(&pivotPort, "port", "8080", "local port to listen on")
	pivotCmd.Flags().StringVar(&pivotProtocol, "protocol", "tcp", "the protocol to use (tcp, udp)")
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
