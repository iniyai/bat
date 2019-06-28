package main

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// General Command interface
type Command interface {

	// Name of the command
	Name() string

	// Single Ling description of this command
	Desc() string

	// Initialize this command
	Init()

	// Interact with STDIN, STDOUT and STDERR
	Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error

	// Print well-descriptive help to STDERR
	Help(stderr io.Writer)
}

// Utility function to run a command with Standard PIPES
func RunCommand(command Command, args []string) chan int {
	retCodeChannel := make(chan int)
	go func() {
		rc := 0
		if command.Interact(args, os.Stdin, os.Stdout, os.Stderr) != nil {
			rc = 1
		}
		retCodeChannel <- rc
	}()
	return retCodeChannel
}

// Common utils -- Start

func ToBuffered(stdin io.Reader, stdout io.Writer, stderr io.Writer) (*bufio.ReadWriter, *bufio.Writer) {
	return bufio.NewReadWriter(bufio.NewReader(stdin), bufio.NewWriter(stdout)), bufio.NewWriter(stderr)
}

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

var CommandNotInitialized = errors.New("Command not initialized")

var IllegalCommandArguments = errors.New("Illegal Command Arguments")

// Common utils -- End
