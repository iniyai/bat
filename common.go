package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Command interface {
	Desc() string

	Init(args []string) error

	Interact(stdin io.Reader, stdout io.Writer, stderr io.Writer) error

	Help(stderr io.Writer)
}

func RunCommand(command Command) chan int {
	retCodeChannel := make(chan int)
	go func() {
		err := command.Interact(os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			fmt.Println("Error while interacting: " + err.Error())
			retCodeChannel <- 1
		} else {
			retCodeChannel <- 0
		}
	}()
	return retCodeChannel
}

// Common utils
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
