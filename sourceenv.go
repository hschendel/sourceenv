package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}
	mustParseDotEnvFile(os.Args[1])
	cmdArgs := os.Args[3:]
	cmd := exec.Command(os.Args[2], cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return
	} else {
		if status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		} else {
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage: sourceenv <.env file> <command> <arg1> .. <argN>")
}

const commentStart = '#'
const envVarRegexpStr = `[_a-zA-Z][_a-zA-Z0-9]*`
var keyValueRegexp = regexp.MustCompile(`^\s*(` + envVarRegexpStr + `)\s*\=\s*(.*)\s*$`) 
var stopLineRegexp = regexp.MustCompile(`^\s*(` + envVarRegexpStr + `)\s*\<\<\<(\S+)\s*$`)

// Reads a .env file and adds its entries to the environment using os.Setenv.
// On error writes to os.Stderr and exits using os.Exit(1).
func mustParseDotEnvFile(filename string) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		fmt.Fprintln(os.Stderr, openErr)
		os.Exit(1)
	}
	reader := bufio.NewReader(file)

	stopLine := ""
	stopLineKey := ""
	stopLineValue := ""
	stopLineStart := 0

	for lineNo := 0; ; lineNo++ {
		line, lineReadErr := reader.ReadString('\n')
		if lineReadErr == nil || lineReadErr == io.EOF {
			if stopLine != "" {
				if strings.HasPrefix(line, stopLine) {
					stopLine = ""
					setenv(stopLineKey, stopLineValue, stopLineStart)
					stopLineKey = ""
					stopLineValue = ""
					stopLineStart = 0
				} else {
					stopLineValue += line
				}
			} else if len(line) > 0 && line[0] != commentStart { // comments in stop line mode are ignored
				matches := keyValueRegexp.FindStringSubmatch(line)
				if matches == nil {
					matches = stopLineRegexp.FindStringSubmatch(line)
					if matches == nil {
						fmt.Fprintln(os.Stderr, "Illformed line", lineNo, ":", line)
						os.Exit(1)
					}
					stopLineKey = matches[1]
					stopLine = matches[2]
					stopLineValue = ""
					stopLineStart = lineNo
				} else {
					setenv(matches[1], matches[2], lineNo)
				}
			}
			if lineReadErr == io.EOF {
				break
			}
		} else {
			fmt.Fprintln(os.Stderr, lineReadErr)
			os.Exit(1)
		}
	}

	file.Close()
}

func setenv(key string, value string, lineNo int) {
	setenvErr := os.Setenv(key, value)
	if setenvErr != nil {
        	fmt.Fprintln(os.Stderr, "Cannot set value from line", lineNo, ":", setenvErr)
                os.Exit(1)
        }
}

