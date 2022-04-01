package search

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/vigilantsys/gogogadget/internal/gadget"
)

var Gadget gadget.Gadget = gadget.Gadget{
	GadgetName:     "search",
	GadgetSynopsis: "search files for a regex pattern",
	GadgetUsage:    "search\n",
	Run:            Run,
	InitFlags:      initFlags,
}

var pattern string
var directory string

//var recursive bool
var ignoreCase bool
var numThreads int

func initFlags(f *flag.FlagSet) {
	f.StringVar(&pattern, "pattern", "", "regex pattern to search for")
	f.StringVar(&directory, "directory", "/", "directory to start search")
	// TODO: Implement non-recursive search
	//f.BoolVar(&recursive, "recursive", false, "whether the search is recursive or not")
	f.BoolVar(&ignoreCase, "ignoreCase", false, "true = case insensitive")
	f.IntVar(&numThreads, "numThreads", 3, "The number of threads to create")
}

func Run() {
	if ignoreCase {
		pattern = fmt.Sprintf("(?i)(%s)", pattern)
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(fmt.Errorf("Error compiling regex expression: %w", err))
		return
	}

	directory = filepath.Clean(directory)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Println(fmt.Errorf("Directory %s does not exist", directory))
	}

	// Create channels
	files := make(chan string)
	searches := make(chan string)

	// Create threads
	var wg sync.WaitGroup

	for w := 1; w <= numThreads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(files, searches, compiled)
		}()
	}

	// Iterate through all files (handle recursiveness)
	go func() {
		defer close(files)
		filepath.WalkDir(directory, func(filepath string, di fs.DirEntry, err error) error {
			if err == nil {
				files <- filepath
			}
			return nil
		})
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
		return "", fmt.Errorf("%w", "File is empty or abnormal")
	}
	if fi.Mode()&fs.ModeSymlink != 0 {
		return "", fmt.Errorf("%w", "File is a symlink")
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
