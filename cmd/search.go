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
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search file contents using regex",
	Long:  `Search the contents of files using a regular expression. Similar to grep.`,
	Run: func(cmd *cobra.Command, args []string) {
		if searchIgnoreCase {
			searchPattern = fmt.Sprintf("(?i)(%s)", searchPattern)
		}
		compiled, err := regexp.Compile(searchPattern)
		if err != nil {
			fmt.Println(fmt.Errorf("error compiling regex expression: %w", err))
			return
		}

		searchDirectory = filepath.Clean(searchDirectory)
		if _, err := os.Stat(searchDirectory); os.IsNotExist(err) {
			fmt.Println(fmt.Errorf("directory %s does not exist", searchDirectory))
		}
		
		// If the path is a file, just do a search on that file
		fileInfo, err := os.Stat(searchDirectory)
		if err != nil {
			fmt.Println(fmt.Errorf("error analyzing path: %w", err))
			return
		}

		if !fileInfo.IsDir() {
			output, err := search(compiled, searchDirectory)
			if err != nil {
				fmt.Println(fmt.Errorf("error searching for pattern in file: %w", err))
				return
			}
			fmt.Println(output)
			return
		}

		// Create channels
		files := make(chan string)
		searches := make(chan string)

		// Create threads
		var wg sync.WaitGroup

		for w := 1; w <= searchThreads; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker(files, searches, compiled)
			}()
		}

		// Iterate through all files
		go func() {
			defer close(files)
			if searchRecursively {
				filepath.WalkDir(searchDirectory, func(filepath string, di fs.DirEntry, err error) error {
					if err == nil {
						files <- filepath
					}
					return nil
				})
			} else {
    				fileInfo, err := ioutil.ReadDir(searchDirectory)
    				if err != nil {
        				fmt.Println(err)
					return
    				}

    				for _, file := range fileInfo {
					if !file.IsDir() {
						fullPath := filepath.Join(searchDirectory, file.Name())
						files <- fullPath
					}
    				}
			}
		}()

		// Wait for threads to complete
		go func() {
			defer close(searches)
			wg.Wait()
		}()

		// Print results
		for search := range searches {
			fmt.Println(search)
			fmt.Println("")
		}
	},
}

var searchPattern, searchDirectory string
var searchIgnoreCase, searchRecursively bool
var searchThreads int

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchPattern, "pattern", "", "pattern to search for")
	searchCmd.Flags().StringVar(&searchDirectory, "directory", ".", "directory to search")
	searchCmd.Flags().BoolVar(&searchIgnoreCase, "ignore-case", false, "ignore case")
	searchCmd.Flags().BoolVar(&searchRecursively, "recursive", false, "Search recursively")
	searchCmd.Flags().IntVar(&searchThreads, "num-threads", runtime.NumCPU(), "number of threads to use")
}

func worker(files <-chan string, searches chan<- string, pattern *regexp.Regexp) {
	for filepath := range files {
		result, err := search(pattern, filepath)
		if err == nil {
			searches <- result
		}
	}
}

func isBinary(line string) bool {
	return strings.Contains(line, "\x00")
}

func search(pattern *regexp.Regexp, filepath string) (string, error) {
	// Prevent read of empty or infinitely long files and symlinks
	fi, err := os.Lstat(filepath)
	if err != nil {
		return "", err
	}
	size := fi.Size()
	if size == 0 {
		return "", fmt.Errorf("%w", "file is empty or abnormal")
	}
	if fi.Mode()&fs.ModeSymlink != 0 {
		return "", fmt.Errorf("%w", "file is a symlink")
	}

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file to find match
	scanner := bufio.NewScanner(file)
	var matches []string
	var lineNum = 1
	for scanner.Scan() {
		line := scanner.Text()
		if isBinary(line) {
			return "", fmt.Errorf("%w", "Binary file")
		}
		if pattern.MatchString(line) {
			matches = append(matches, fmt.Sprintf("%s:%d:%s\n", filepath, lineNum, scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		lineNum++
	}

	// Return matches
	if len(matches) > 0 {
		output := ""
		for _, match := range matches {
			output += match
		}
		return output, nil
	}

	// No matches found
	return "", fmt.Errorf("%w", "No match")
}
