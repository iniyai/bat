package main

import (
	"bufio"
	"io"
	"os"
)

// General Command interface
type Command interface {
	// Single Ling description of this command
	Desc() string

	// Initialize this command with cmd line arguments
	Init(args []string) error

	// Interact with STDIN, STDOUT and STDERR
	Interact(stdin io.Reader, stdout io.Writer, stderr io.Writer) error

	// Print well-descriptive help to STDERR
	Help(stderr io.Writer)
}

// Utility function to run a command with Standard PIPES
func RunCommand(command Command) chan int {
	retCodeChannel := make(chan int)
	go func() {
		rc := 0
		if command.Interact(os.Stdin, os.Stdout, os.Stderr) != nil {
			rc = 1
		}
		retCodeChannel <- rc
	}()
	return retCodeChannel
}

// Common utils -- Start
func WriteAndFlush(writer *bufio.Writer, value string) {
	_, err := writer.WriteString(value)

	if err != nil {
		panic("I/O error")
	} else {
		err := writer.Flush()
		if err != nil {
			panic("I/O error")
		}
	}
}
