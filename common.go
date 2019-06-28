package main

import (
	"bufio"
	"errors"
	"github.com/fatih/color"
	"io"
)

// Common utils -- Start


// Output Formatting options -- Start
var BlackBoldFormatter = color.New(color.FgBlack).Add(color.Bold)
var ItalicFormatter = color.New(color.Italic)
var BlackBoldUnderLineFormatter = color.New(color.FgBlack).Add(color.Bold).Add(color.Underline)
// Output Formatting options -- End

// IOUtils -- Start
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
// IOUtils -- End

// Common Errors -- Start
var CommandNotInitialized = errors.New("Command not initialized")

var IllegalCommandArguments = errors.New("Illegal Command Arguments")
// Common Errors -- End

// Common utils -- End
