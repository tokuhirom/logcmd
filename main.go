package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type stripWriter struct {
	w io.Writer
}

func (s *stripWriter) Write(p []byte) (n int, err error) {
	cleaned := ansiEscape.ReplaceAll(p, nil)
	_, err = s.w.Write(cleaned)
	return len(p), err
}

func main() {
	// Manual flag parsing (flag package doesn't play well with subcommands)
	noHeader := false
	noExit := false

	args := os.Args[1:]
	filtered := args[:0]
	for _, a := range args {
		switch a {
		case "--no-header":
			noHeader = true
		case "--no-exit":
			noExit = true
		default:
			filtered = append(filtered, a)
		}
	}

	if len(filtered) < 2 {
		fmt.Fprintln(os.Stderr, "usage: logcmd [--no-header] [--no-exit] <logfile> <command> [args...]")
		os.Exit(1)
	}

	logPath := filtered[0]
	cmdArgs := filtered[1:]

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	sw := &stripWriter{w: f}

	if !noHeader {
		if fi, err := f.Stat(); err == nil && fi.Size() > 0 {
			fmt.Fprintln(sw, "")
		}

		fmt.Fprintf(sw, "=== [%s] $ %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(cmdArgs, " "))
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, sw)
	cmd.Stderr = io.MultiWriter(os.Stderr, sw)
	cmd.Stdin = os.Stdin

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 1
		}
	}

	if !noExit {
		fmt.Fprintf(sw, "=== exit code: %d\n", exitCode)
	}

	os.Exit(exitCode)
}
