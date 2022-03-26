package search 


import (
	"flag"
	"fmt"
	"path/filepath"
	"io/fs"
	"sync"
	"regexp"
	"bufio"
	"os"
	"io/ioutil"
	"bytes"

	"github.com/seandheath/gogogadget/internal/gadget"
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
		filepath.WalkDir(directory, func(filepath string, di fs.DirEntry, err error) error {
			if err == nil {
        			files <- filepath
    			}
			return nil
		})
		close(files)
	}()

	// Wait for threads to complete
	go func() {
		wg.Wait()
		close(searches)
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

func isBinary(file string) (bool, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}
	return bytes.Contains(b, []byte("\x00")), nil
}


func search(pattern *regexp.Regexp, filepath string) (string, error) {
	if isB, err := isBinary(filepath); isB {
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("%w", "File is a binary file")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)
	matches := make(map[int]string)
	var lineNum = 1
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			matches[lineNum] = scanner.Text()
		} 
		if err := scanner.Err(); err != nil {
			return "", err
		}
		lineNum++
	}
	if len(matches) > 0 {
		output := filepath + ":\n"
		for ln, line := range matches {
			output += fmt.Sprintf("%d:%s\n", ln, line)
		}
		return output, nil
	}
	return "", fmt.Errorf("%w", "No match")
}
