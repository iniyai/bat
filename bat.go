package main

import (
	"fmt"
	"os"
	"strings"
)

var WelcomeMessage = "Welcome to " + BlackBoldUnderLineFormatter.Sprint("BAT") +
	" (" +
	ItalicFormatter.Sprint("Bash Additional Tools") +
	").\n\n" +
	"Give me one of below commands.\n" +
	"For more help type bat help <cmd>"

// To enable your command, add an entry here.
var enabledCommands = []Command{
	&StatCommand{},
	&LinesBetweenCommand{},
	&HistogramCommand{},
	&LineLengthCommand{},
	&EnvironmentCommand{},
}

// Initialize Enabled commands
func initCommands() map[string]Command {
	commands := make(map[string]Command)

	for _, cmd := range enabledCommands {
		cmd.Init()
		commands[cmd.Name()] = cmd
	}

	return commands
}

func main() {

	// All commands
	cmds := initCommands()

	// If no commands are given show general help.
	if len(os.Args) < 2 {
		fmt.Println(WelcomeMessage)
		index := 0
		for name, c := range cmds {
			index++
			fmt.Printf("  %d) %s - %s\n", index, BlackBoldFormatter.Sprint(name), c.Desc())
		}
		os.Exit(1)
	}

	cmdName := strings.ToLower(os.Args[1])
	// Handle given command
	switch cmdName {

	default:
		cmd, ok := cmds[cmdName]
		if ok {
			os.Exit(<-RunCommand(cmd, os.Args[2:]))
		} else {
			fmt.Println("unknown command: " + cmdName)
		}
		break

	case "help":
		if len(os.Args) == 2 {
			fmt.Println("specify a command for help\n bat help <cmd>")
		} else {
			helpForCmd := os.Args[2]
			cmd, ok := cmds[helpForCmd]
			if ok {
				cmd.Help(os.Stderr)
			} else {
				fmt.Println("unknown command: " + helpForCmd)
			}
		}
	}
}
