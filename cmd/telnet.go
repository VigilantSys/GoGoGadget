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
	//"bufio"
	"os"
	"strings"
	//"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/reiver/go-telnet"
	"github.com/reiver/go-oi"
)


// telnetCmd represents the telnet command
var telnetCmd = &cobra.Command{
	Use:   "telnet",
	Short: "telnet client",
	Long:  `Telnet provides a telnet client program.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := telnet.DialTo(fmt.Sprintf("%s:%s", telnetAddress, telnetPort))
		if err != nil {
			fmt.Println(fmt.Errorf("Error connecting to telnet server: %w", err))
			return
		}

		go Receive(conn)
		c := make(chan bool, 1)
		go Send(conn, c)
		<-c

	},
}

var telnetAddress string
var telnetPort string

func init() {
	rootCmd.AddCommand(telnetCmd)
	telnetCmd.Flags().StringVar(&telnetAddress, "address", "", "address to connect to")
	telnetCmd.Flags().StringVar(&telnetPort, "port", "23", "port to connect to")
}

func Receive(conn *telnet.Conn) {
    //var buffer [1]byte
    //recvData := buffer[:]
    recvData := make([]byte, 1)
    var n int
    var err error

    for {
        n, err = conn.Read(recvData)
        //fmt.Println("Bytes: ", n, "Data: ", recvData, string(recvData))
        if n <= 0 || err != nil {
            break
        } else {
		//r, _ :=utf8.DecodeRune(recvData)
		//os.Stdout.Write(recvData)
		if _, err := oi.LongWriteString(os.Stdout, string(recvData)); nil != err {
		fmt.Println(err)
	}
        }
    }
}

func Send(conn *telnet.Conn, c chan bool) {
	//reader := bufio.NewReader(os.Stdin)
	defer conn.Close()
	running := ""
	for {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		if string(b) == "\n" || string(b) == "\r" || string(b) == "\r\n" {
			if strings.TrimSpace(running) == "quit" || strings.TrimSpace(running) == "exit" {
				break
			}
			running = ""
		}
		running += string(b)

		conn.Write(b)
		/*
		cmd, _ := reader.ReadString('\n')
		if strings.TrimSpace(cmd) == "exit" || strings.TrimSpace(cmd) == "quit" {
			break
		}
		*/
    
		/*
		var commandBuffer []byte
    		for _, char := range cmd {
        		commandBuffer = append(commandBuffer, byte(char))
    		}

    		var crlfBuffer [2]byte = [2]byte{'\r', '\n'}
    		crlf := crlfBuffer[:]

    		conn.Write(commandBuffer)
    		conn.Write(crlf)
		*/
		//conn.Write([]byte(cmd + "\r\n"))
    	}
	c <- true
}
