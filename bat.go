package main

import (
	"fmt"
	"os"
	"strings"
)

const WelcomeMessage string = "Welcome to BAT (Bash Additional Tools).\n\n" +
	"Give me one of below commands.\n" +
	"For more help type bat help <cmd>"

func buildCmds() map[string]Command {
	commands := make(map[string]Command)

	for _, cmd := range []Command{&StatCommand{},
		&LinesBetweenCommand{}} {
		cmd.Init()
		commands[cmd.Name()] = cmd
	}

	return commands
}

func main() {

	// All commands
	cmds := buildCmds()

	// If no commands are given show general help.
	if len(os.Args) < 2 {
		fmt.Println(WelcomeMessage)
		for name, c := range cmds {
			fmt.Println("  " + name + " - " + c.Desc())
		}
		os.Exit(1)
	}

	cmd_name := strings.ToLower(os.Args[1])
	// Handle given command
	switch cmd_name {

	default:
		cmd, ok := cmds[cmd_name]
		if ok {
			os.Exit(<-RunCommand(cmd, os.Args[2:]))
		} else {
			fmt.Println("unknown command: " + cmd_name)
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
