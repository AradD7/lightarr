package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *config, param ...string) error {
	fmt.Println("Closing Lightarr... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, param ...string) error {
	fmt.Printf("Welcome to the Lightarr!\nUsage:\n\n")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}

func commandPrintBulbs(cfg *config, param ...string) error {
	for _, bulb := range cfg.bulbsMap {
		fmt.Println(*bulb)
	}
	return nil
}
