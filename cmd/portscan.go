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
	"net"
	"os"
	//"sort"
	"errors"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var portscanCmd = &cobra.Command{
	Use:   "portscan",
	Short: "A simple TCP scanner",
	Long:  `A scanner for TCP ports. Similar to nmap.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Find good value for capacity
		// Initialize variables
		const capacity = 100
		portsList, err := parsePorts(portscanPorts)
		if err != nil {
			fmt.Println(fmt.Errorf("Invalid port string: %w", err))
			return
		}
		ipsList, err := parseIPs(portscanAddress)
		if err != nil {
			fmt.Println(fmt.Errorf("Invalid IP range: %w", err))
			return
		}
		timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ds", portscanTimeout))
		if err != nil {
			fmt.Println(fmt.Errorf("Timeout not valid: %w", err))
			return
		}

		// Perform ping scan
		if !portscanNoPing {
			fmt.Printf("Performing ping scan on %d IP addresses\n\n", len(ipsList))
			ipsList, err = pingScan(ipsList, capacity, timeoutDuration)
			if err != nil {
				fmt.Println(fmt.Errorf("Error performing ping scan: %w", err))
				return
			}
			fmt.Printf("Found %d hosts up\n\n", len(ipsList))
			fmt.Println("Ping scan complete!\n")
		}

		// Perform port scan
		fmt.Printf("Performing portscan for %d ports on %d IP addresses:\n\n", len(portsList), len(ipsList))
		resultMap, err := portScan(ipsList, portsList, capacity, timeoutDuration)
		if err != nil {
			fmt.Println(fmt.Errorf("Error performing port scan: %w", err))
			return
		}
		fmt.Println("\nPort scan complete!\n")

		// Print Results
		fmt.Println("Scan results:")
		printResults(resultMap)
	},
}

var portscanAddress, portscanPorts string
var portscanTimeout int
var portscanNoPing bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogogadget.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(portscanCmd)
	portscanCmd.Flags().StringVar(&portscanAddress, "ip", "", "The IP address or range to scan (CIDR).")
	portscanCmd.Flags().StringVar(&portscanPorts, "ports", "", "The range of ports to scan.")
	portscanCmd.Flags().IntVar(&portscanTimeout, "timeout", 1, "The timeout length in seconds.")
	portscanCmd.Flags().BoolVar(&portscanNoPing, "noping", false, "Whether to do a pingscan first")
}

func portscanWorker(addresses, results chan string, timeout time.Duration) {
	for address := range addresses {
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				results <- fmt.Sprintf("%s %s", address, "timeout")
				continue
			}
			results <- fmt.Sprintf("%s %s", address, "closed")
			continue
		}
		conn.Close()
		results <- fmt.Sprintf("%s %s", address, "open")
		fmt.Printf("Found open port: %s\n", address)
	}
}

func portScan(ipsList []string, portsList []string, capacity int, timeout time.Duration) (map[string]map[string]string, error) {
	addresses := make(chan string, capacity)
	results := make(chan string)
	defer close(addresses)
	defer close(results)
	resultMap := make(map[string]map[string]string)

	for i := 0; i < cap(addresses); i++ {
		go portscanWorker(addresses, results, timeout)
	}

	go func() {
		for i := range ipsList {
			for p := range portsList {
				address := fmt.Sprintf("%s:%s", ipsList[i], portsList[p])
				addresses <- address
			}
		}
	}()

	for i := 0; i < len(portsList)*len(ipsList); i++ {
		result := <-results
		if result != "" {
			fields := strings.Split(result, " ")
			portStatus := fields[1]
			fields = strings.Split(fields[0], ":")
			// TODO: add error handling to this
			if resultMap[fields[0]] == nil {
				resultMap[fields[0]] = make(map[string]string)
			}
			resultMap[fields[0]][fields[1]] = portStatus
		}
	}

	return resultMap, nil
}

func pingScan(ips []string, capacity int, timeout time.Duration) ([]string, error) {
	ports := []string{"80", "443"}
	resultMap, err := portScan(ips, ports, capacity, timeout)
	if err != nil {
		return nil, err
	}

	var aliveIps []string
	for ip, portMap := range resultMap {
		for _, status := range portMap {
			if status != "timeout" {
				aliveIps = append(aliveIps, ip)
				break
			}
		}
	}

	return aliveIps, nil
}

func parsePorts(portsString string) ([]string, error) {
	/* Credit: https://github.com/blackhat-go/bhg/blob/master/ch-2/scanner-port-format/portformat.go */
	ports := []string{}
	if strings.Contains(portsString, ",") && strings.Contains(portsString, "-") {
		sp := strings.Split(portsString, ",")
		for _, p := range sp {
			if strings.Contains(p, "-") {
				if err := dashSplit(p, &ports); err != nil {
					return ports, err
				}
			} else {
				if err := convertAndAddPort(p, &ports); err != nil {
					return ports, err
				}
			}
		}
	} else if strings.Contains(portsString, ",") {
		sp := strings.Split(portsString, ",")
		for _, p := range sp {
			convertAndAddPort(p, &ports)
		}
	} else if strings.Contains(portsString, "-") {
		if err := dashSplit(portsString, &ports); err != nil {
			return ports, err
		}
	} else {
		if err := convertAndAddPort(portsString, &ports); err != nil {
			return ports, err
		}
	}
	return ports, nil

}

func parseIPs(ipsString string) ([]string, error) {
	/* Credit : https://gist.github.com/kotakanbe/d3059af990252ba89a82 */

	if !strings.Contains(ipsString, "/") {
		ipsString += "/32"
	}
	ip, ipnet, err := net.ParseCIDR(ipsString)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	return ips, nil
}

func inc(ip net.IP) {
	/* Credit : https://gist.github.com/kotakanbe/d3059af990252ba89a82 */
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func printResults(results map[string]map[string]string) {
	fmt.Println("")

	for ip, portMap := range results {
		fmt.Printf("Results for host %s:\n", ip)
		var portStrings []string
		for port, status := range portMap {
			if status != "closed" {
				portStrings = append(portStrings, fmt.Sprintf("%s\t%s\t", port, status))
			}
		}

		// If all ports are closed, don't print table
		if len(portStrings) == 0 {
			fmt.Printf("All %d scanned ports are closed\n\n", len(portMap))
			continue
		}

		// Print port table
		const padding = 5
		table := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		fmt.Fprintln(table, "PORT\tSTATUS\t")
		for line := range portStrings {
			fmt.Fprintln(table, portStrings[line])
		}
		table.Flush()
		fmt.Println("\n")
	}
}

/* Credit: https://github.com/blackhat-go/bhg/blob/master/ch-2/scanner-port-format/portformat.go */
const (
	porterrmsg = "Invalid port specification"
)

func dashSplit(sp string, ports *[]string) error {
	dp := strings.Split(sp, "-")
	if len(dp) != 2 {
		return errors.New(porterrmsg)
	}
	start, err := strconv.Atoi(dp[0])
	if err != nil {
		return errors.New(porterrmsg)
	}
	end, err := strconv.Atoi(dp[1])
	if err != nil {
		return errors.New(porterrmsg)
	}
	if start > end || start < 1 || end > 65535 {
		return errors.New(porterrmsg)
	}
	for ; start <= end; start++ {
		*ports = append(*ports, strconv.Itoa(start))
	}
	return nil
}

func convertAndAddPort(p string, ports *[]string) error {
	i, err := strconv.Atoi(p)
	if err != nil {
		return errors.New(porterrmsg)
	}
	if i < 1 || i > 65535 {
		return errors.New(porterrmsg)
	}
	*ports = append(*ports, strconv.Itoa(i))
	return nil
}
