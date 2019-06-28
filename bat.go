package main

import (
	"fmt"
	"os"
	"strings"
)

func buildCmds() map[string]Command {
	commands := make(map[string]Command)
	commands["stat"] = &StatCommand{}
	return commands
}

func main() {

	cmds := buildCmds()

	if len(os.Args) < 2 {
		fmt.Println("Welcome to BAT (Bash Additional Tools).\n\nGive me one of below commands.\nFor more help type bat help <cmd>")
		for name, c := range cmds {
			fmt.Println("  " + name + " - " + c.Desc())
		}
		os.Exit(1)
	}

	cmd_name := strings.ToLower(os.Args[1])
	switch cmd_name {

	default:
		cmd, ok := cmds[cmd_name]
		if !ok {
			fmt.Println("unknown command: " + cmd_name)
		} else {
			err := cmd.Init(os.Args[2:])
			if err != nil {
				fmt.Println("unable to initialize cmd: " + cmd_name + " with args: " + strings.Join(os.Args[2:], ","))
			} else {
				os.Exit(<- RunCommand(cmd))
			}
			cmd.Help(os.Stderr)
			break
		}

	case "help":

	}

}
