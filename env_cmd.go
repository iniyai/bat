package main

import (
	"flag"
	"fmt"
	"io"
	"runtime"
)

type EnvironmentCommand struct {
	cmdOptions *flag.FlagSet
}

func (EnvironmentCommand) Name() string {
	return "env"
}

func (EnvironmentCommand) Desc() string {
	return "Information about current running environment"
}

func (ec *EnvironmentCommand) Init() {
	ec.cmdOptions = flag.NewFlagSet(ec.Name(), flag.ExitOnError)
}

func (ec *EnvironmentCommand) Interact(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if ec.cmdOptions.Parse(args) != nil {
		return CommandNotInitialized
	}

	bInOut, _ := ToBuffered(stdin, stdout, stderr)
	WriteAndFlush(bInOut.Writer, fmt.Sprintf("%s running on %s\n", BlackBoldFormatter.Sprint(runtime.GOOS),
		BlackBoldFormatter.Sprint(runtime.GOARCH)))

	WriteAndFlush(bInOut.Writer, fmt.Sprintf("Number of cores: %s\n",
		BlackBoldFormatter.Sprint(runtime.NumCPU())))

	return nil
}

func (ec *EnvironmentCommand) Help(stderr io.Writer) {
	ec.cmdOptions.SetOutput(stderr)
	ec.cmdOptions.Usage()
}
