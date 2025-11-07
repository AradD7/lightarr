package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/AradD7/lightarr/internal/wiz"
)


type cliCommand struct {
	name			string
	description		string
	callback		func(cfg *config, param ...string) error
}

type config struct {
	conn 		*net.UDPConn
	bulbsMap	map[string]*wiz.Bulb
}

func cleanInput(text string) []string {
	seperatedText := strings.Split(strings.Trim(strings.ToLower(text), " "), " ")
	seperatedText = slices.DeleteFunc(seperatedText, func(s string) bool {
		return s == ""
	})
	return seperatedText
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("lightarr > ")

		scanner.Scan()

		userInputSlice := cleanInput(scanner.Text())
		if len(userInputSlice) == 0 {
			continue
		}


		if c, ok := getCommands()[userInputSlice[0]]; ok {
			fmt.Println()
			err := c.callback(cfg, userInputSlice...)
			fmt.Println()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command")
			}
	}
}


func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"exit": {
			name: "exit",
			description: "Exit Lightarr",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays help message",
			callback: commandHelp,
		},
		"bulbs": {
			name: "bulbs",
			description: "Prints all the bulbs",
			callback: commandPrintBulbs,
		},
	}
}

